package utils

import "fmt"

// func Filter(s []interface{}, f func(interface{})bool)(r []interface{}){

// 	for _, o := range s {

// 		if f(o){
// 			r = append(r, o)
// 		}

// 	}

// 	return

// }

type i interface{}

func Filter(s []interface{}) {

	for _, o := range s {

		fmt.Println(o)

	}

}
