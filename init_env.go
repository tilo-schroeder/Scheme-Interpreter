package main

var globalEnv env

func init() {
	globalEnv = env{
		vars{
			"+": func(a ...expr) expr {
				v := a[0].(float64)
				for _, i := range a[1:] {
					v += i.(float64)
				}
				return v
			},
			"-": func(a ...expr) expr {
				v := a[0].(float64)
				for _, i := range a[1:] {
					v -= i.(float64)
				}
				return v
			},
			"*": func(a ...expr) expr {
				v := a[0].(float64)
				for _, i := range a[1:] {
					v *= i.(float64)
				}
				return v
			},
			"/": func(a ...expr) expr {
				v := a[0].(float64)
				for _, i := range a[1:] {
					v /= i.(float64)
				}
				return v
			},
			//What to do with more than two parameters?
			"<": func(a ...expr) expr {
				return Bool(a[0].(float64) < a[1].(float64))
			},
			">": func(a ...expr) expr {
				return Bool(a[0].(float64) > a[1].(float64))
			},
			"<=": func(a ...expr) expr {
				return Bool(a[0].(float64) <= a[1].(float64))
			},
			">=": func(a ...expr) expr {
				return Bool(a[0].(float64) >= a[1].(float64))
			},
			"=": func(a ...expr) expr {
				return Bool(a[0].(float64) == a[1].(float64))
			},
			"not": func(a ...expr) expr {
				return !a[0].(Bool)
			},
			/*"and": func(a ...expr) expr {
				res := Bool(true)
				for i := range a {
					res = res && a[i].(Bool)
				}
				return res
			},*/
			"and": func(a ...expr) expr {
				for i := range a {
					if a[i].(Bool) == Bool(false) {
						return Bool(false)
					}
				}
				return Bool(true)
			},
			/*"or": func(a ...expr) expr {
				res := Bool(false)
				for i := range a {
					res = res || a[i].(Bool)
				}
				return res
			},*/
			"or": func(a ...expr) expr {
				for i := range a {
					if a[i].(Bool) == Bool(true) {
						return Bool(true)
					}
				}
				return Bool(false)
			},
			"symbol?": func(a ...expr) expr {
				switch a[0].(type) {
				case Symbol:
					return true
				default:
					return false
				}
			},
			"first": func(a ...expr) expr {
				return a[0].([]expr)[0]
			},
			"rest": func(a ...expr) expr {
				return a[0].([]expr)[1:]
			},
			"cons": func(a ...expr) expr {
				switch first := a[0]; rest := a[1].(type) {
				case []expr:
					return append([]expr{first}, rest...)
				default:
					return []expr{first, rest}
				}
			},
			"list": func(a ...expr) expr {
				list := []expr{}
				list = append(list, a...)
				return list
			},
		},
		nil}
}
