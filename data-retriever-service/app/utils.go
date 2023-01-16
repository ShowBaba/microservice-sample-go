package app

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/graphql-go/graphql"
	"gorm.io/gorm"
)

func makeNodeListType(name string, nodeType *graphql.Object) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: name,
		Fields: graphql.Fields{
			"nodes":      &graphql.Field{Type: graphql.NewList(nodeType)},
			"totalCount": &graphql.Field{Type: graphql.Int},
		},
	})
}

func makeListField(listType graphql.Output, resolve graphql.FieldResolveFn) *graphql.Field {
	var (
		listField = &graphql.Field{
			Type:    listType,
			Resolve: resolve,
			Args: graphql.FieldConfigArgument{
				"limit":   &graphql.ArgumentConfig{Type: graphql.Int},
				"offset":  &graphql.ArgumentConfig{Type: graphql.Int},
				"orderBy": &graphql.ArgumentConfig{Type: graphql.String},
			},
		}
		fields graphql.FieldDefinitionMap
	)

	switch listType.Name() {
	case "UserList":
		fields = UserType.Fields()
	case "PostList":
		fields = PostType.Fields()
	}

	for key, val := range fields {
		if _, ok := listField.Args[key]; !ok {
			listField.Args[key] = &graphql.ArgumentConfig{Type: val.Type}
		}
	}
	return listField
}

func parseDbClause(params graphql.ResolveParams, tx *gorm.DB, nodeType *graphql.Object) *gorm.DB {
	// set all query fields
	if limit, ok := params.Args["limit"].(int); ok {
		tx = tx.Limit(limit)
	}

	if offset, ok := params.Args["offset"].(int); ok {
		tx = tx.Offset(offset)
	}

	if orderBy, ok := params.Args["orderBy"].(string); ok {
		direction := "ASC"
		if strings.Index(orderBy, "-") == 0 {
			direction = "DESC"
			orderBy = orderBy[1:]
		}
		tx = tx.Order(orderBy + " " + direction)
	}

	var fields graphql.FieldDefinitionMap = nodeType.Fields()

	for key, field := range fields {
		val, ok := params.Args[key]
		if !ok {
			continue
		}
		if key == "firstname" && nodeType.Name() == "User" {
			tx = tx.Where("LOWER(firstname) LIKE ?", fmt.Sprintf(`%%%s%%`, strings.ToLower(val.(string))))
			continue
		}
		if key == "lastname" && nodeType.Name() == "User" {
			tx = tx.Where("LOWER(lastname) LIKE ?", fmt.Sprintf(`%%%s%%`, strings.ToLower(val.(string))))
			continue
		}
		if key == "title" && nodeType.Name() == "Post" {
			tx = tx.Where("LOWER(title) LIKE ?", fmt.Sprintf(`%%%s%%`, strings.ToLower(val.(string))))
			continue
		}
		key = underscore(key)
		switch field.Type {
		case graphql.String:
			if match := regexp.MustCompile(`((>|<)=?)\s*(.*?)$`).FindStringSubmatch(val.(string)); len(match) > 0 {
				tx = tx.Where(key+" "+match[1]+"?", match[3])
			} else {
				tx = tx.Where(key+" = ?", val.(string))
			}
		case graphql.Boolean:
			tx = tx.Where(key+" = ?", val.(bool))
		case graphql.Int:
			// handle >= or <= int value
			if nodeType.Name() == "MintCalendarCollectionList" {
				break
			}

			if intVal, ok := val.(int); ok {
				tx = tx.Where(key+" = ?", intVal)
			} else if strVal, ok := val.(string); ok {
				if intVal, err := strconv.Atoi(strVal); err == nil {
					tx = tx.Where(key+" = ?", intVal)
				} else if match := regexp.MustCompile(`(>=?)\s*(\d+),(<=?)\s*(\d+)`).FindStringSubmatch(strVal); len(match) > 0 {
					intVal, _ = strconv.Atoi(match[2])
					tx = tx.Where(key+" "+match[1]+" ?", intVal)

					intVal, _ = strconv.Atoi(match[4])
					tx = tx.Where(key+" "+match[3]+" ?", intVal)
				} else if match := regexp.MustCompile(`((>|<)=?)\s*(\d+)`).FindStringSubmatch(strVal); len(match) > 0 {
					intVal, _ = strconv.Atoi(match[3])
					tx = tx.Where(key+" "+match[1]+" ?", intVal)
				}
			}
		case graphql.Float:
			// handle >= or <= int value
			if floatVal, ok := val.(float64); ok {
				tx = tx.Where(key+" = ?", floatVal)
			} else if strVal, ok := val.(string); ok {
				if floatVal, err := strconv.ParseFloat(strVal, 64); err == nil {
					tx = tx.Where(key+" = ?", floatVal)
				} else if match := regexp.MustCompile(`(>=?)\s*(\d+(\.\d+)?),(<=?)\s*(\d+(\.\d+)?)`).FindStringSubmatch(strVal); len(match) > 0 {
					floatVal, _ = strconv.ParseFloat(match[2], 64)
					tx = tx.Where(key+" "+match[1]+" ?", floatVal)

					floatVal, _ = strconv.ParseFloat(match[5], 64)
					tx = tx.Where(key+" "+match[4]+" ?", floatVal)
				} else if match := regexp.MustCompile(`((>|<)=?)\s*(\d+(\.\d+)?)`).FindStringSubmatch(strVal); len(match) > 0 {
					floatVal, _ = strconv.ParseFloat(match[3], 64)
					tx = tx.Where(key+" "+match[1]+" ?", floatVal)
				}
			}
		}
	}

	return tx
}

var camel = regexp.MustCompile("(^[^A-Z0-9]*|[A-Z0-9]*)([A-Z0-9][^A-Z]+|$)")

func underscore(s string) string {
	var a []string
	for _, sub := range camel.FindAllStringSubmatch(s, -1) {
		if sub[1] != "" {
			a = append(a, sub[1])
		}
		if sub[2] != "" {
			a = append(a, sub[2])
		}
	}
	return strings.ToLower(strings.Join(a, "_"))
}
