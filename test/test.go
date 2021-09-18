package main

import (
	"context"
	"fmt"

	elastic7 "github.com/olivere/elastic/v7"
)

type Meta map[string]interface{}

func main() {
	client, err := elastic7.NewClient(elastic7.SetURL("http://192.168.31.21:9200"), elastic7.SetSniff(false))
	if err != nil {
		fmt.Println(err)
	}
	_, results, _ := aliasElasticsearchGet(client, "logs_write")
	fmt.Println(results)
}

func aliasElasticsearchGet(client *elastic7.Client, name string) (bool, []Meta, error) {
	fmt.Println(name)
	var aliasMeta []Meta
	result, err := client.CatAliases().Alias(name).Do(context.TODO())
	if err != nil {
		return false, nil, err
	}
	for _, data := range result {
		res := make(map[string]interface{})
		res["alias"] = data.Alias
		res["index"] = data.Index
		res["is_write_index"] = data.IsWriteIndex
		res["routing.index"] = data.RoutingIndex
		res["routing.search"] = data.RoutingSearch
		aliasMeta = append(aliasMeta, res)
	}
	return len(result) > 0, aliasMeta, nil
}
