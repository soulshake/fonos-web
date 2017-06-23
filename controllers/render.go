package controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/auroralaboratories/pulse"
)

type Data struct {
	Body                string
	Template            string
	Error               error
	FlashSuccessMessage string
	FlashErrorMessage   string
	Input               string
	Output              string
	Sinks               Sinks
	Volumes             []Volume
	SinkInputs          []pulse.Source
	//JSON                string
}

func RenderHTML(w http.ResponseWriter, r *http.Request, layout string, templateName string, data Data) (err error) {
	//log.Printf("*** Rendering templateName: %s", templateName)
	t := template.New(templateName).Funcs(templateHelpers())

	t, err = t.ParseFiles(
		//fmt.Sprintf("%s/layouts/%s.html", TemplateDir, layout),
		fmt.Sprintf("%s/header.html", TemplateDir),
		fmt.Sprintf("%s/dashboard.html", TemplateDir),
		fmt.Sprintf("%s/api.html", TemplateDir),
		fmt.Sprintf("%s/index.html", TemplateDir),
		fmt.Sprintf("%s/volumes.html", TemplateDir),
		fmt.Sprintf("%s/sinks.html", TemplateDir),
		fmt.Sprintf("%s/sink-inputs.html", TemplateDir),
		fmt.Sprintf("%s/pacmd.html", TemplateDir),
		fmt.Sprintf("%s/flash-error.html", TemplateDir),
		fmt.Sprintf("%s/flash-error-auth.html", TemplateDir),
		fmt.Sprintf("%s/"+templateName+".html", TemplateDir),
	)

	if err != nil {
		log.Fatal(err)
	}

	err = t.ExecuteTemplate(w, fmt.Sprintf("%s", layout), data)
	if err != nil {
		log.Print(err.Error())
		return err
	}

	return nil
}

func RenderStatic(w http.ResponseWriter, r *http.Request, templateName string) error {
	data := Data{Template: templateName}
	//data := map[string]interface{}{
	//"Template": templateName,
	//}

	return RenderHTML(w, r, "dashboard", templateName, data)
}

func templateHelpers() template.FuncMap {
	return template.FuncMap{
		"env": func(key string) string {
			return os.Getenv(key)
		},
		/*
			"date": func(s int64) string {
				t := time.Unix(s, 0)
				return t.Format("Jan 2, 2006")
			},
			"dateISO8601": func(t time.Time) string {
				return t.Format("2006-01-02T15:04:05Z")
			},
		*/
		"defined": func(data interface{}, key string) bool {
			if params, ok := data.(map[string]interface{}); ok {
				if _, ok := params[key]; ok {
					return true
				}
			}
			return false
		},
		/*
			"discount": func(subtotal, total int64) string {
				cents := subtotal - total
				d := cents / 100
				c := cents - d*100
				return fmt.Sprintf("-$%d.%02d", d, c)
			},
			"dollars": func(cents uint64) string {
				d := cents / 100
				c := cents - d*100
				return fmt.Sprintf("$%d.%02d", d, c)
			},
			"udollars": func(cents uint64) string {
				d := cents / 100
				c := cents - d*100
				return fmt.Sprintf("$%d.%02d", d, c)
			},
			"duration": func(start, end time.Time) string {
				if end.IsZero() {
					return ""
				}

				s := int(end.Sub(start).Seconds())

				if s > 60 {
					return fmt.Sprintf("%dm %02ds", s/60, s%60)
				}

				return fmt.Sprintf("%ds", s)
			},
			"md5": func(s string) string {
				return fmt.Sprintf("%x", md5.Sum([]byte(s)))
			},
		*/
		"upper": func(t interface{}) string {
			s := fmt.Sprintf("%s", t) // convert any passed in Type to a string representation
			return strings.ToUpper(s[0:1]) + s[1:]
		},
		"floatToPrecision": func(f float64, places int) string {
			fmtStr := fmt.Sprintf(`%%.%df`, places)
			return fmt.Sprintf(fmtStr, f)
		},
		"multipliedBy": func(f, by float64) float64 {
			return f * by
		},
		"add": func(x, y int) int {
			return x + y
		},
		/*
			"role": func(role string, admin bool) template.HTML {
				html := `<em>`
				html += `<span class=" label label-role label-role-` + role + `"> `
				html += `<span class="role-name">` + role + `</span>`
				html += `</span>`
				if admin {
					html += ` <a class="toggle-update-role">(edit)</a>`
				}
				html += `</em>`
				return template.HTML(html)
			},
		*/
	}
}

var knownPaths = map[string]string{
	"POST:/apps":           "Created application <b>{name}</b>",
	"DELETE:/apps/{id}":    "Deleted application <b>{id}</b>",
	"GET:/apps/{app}/logs": "Ran <code>logs</code> for <b>{app}</b>",
	"PUT:/system":          "Updated rack <b>{params}</b>",
}

/*
func matchPath(log models.AuditLog) string {
	sig := fmt.Sprintf("%s:%s", log.Method, log.Path)

	body, err := url.ParseQuery(log.Body)
	if err != nil {
		return ""
	}

	if strings.HasSuffix(log.Path, "/environment") {
		env := strings.Split(strings.TrimSpace(log.Body), "\n")
		keys := []string{}
		for _, e := range env {
			keys = append(keys, strings.SplitN(e, "=", 2)[0])
		}
		body.Add("vars", strings.Join(keys, ", "))
	} else {
		disp := []string{}
		for k := range body {
			disp = append(disp, fmt.Sprintf("%s=%s", k, body.Get(k)))
		}
		body.Add("params", strings.Join(disp, " "))
	}

	for path, desc := range knownPaths {
		re, tokens := pathRegexp(path)
		if re == nil {
			continue
		}

		m := re.FindAllStringSubmatch(sig, -1)
		if len(m) == 0 {
			continue
		}

		if len(m[0]) > 1 {
			for i, v := range m[0][1:] {
				body.Add(tokens[i], v)
			}
		}

		desc = strings.Replace(desc, "{body}", log.Body, -1)

		return regPathTokenParser.ReplaceAllStringFunc(desc, func(s string) string {
			if len(s) > 2 {
				return body.Get(s[1 : len(s)-1])
			}
			return ""
		})
	}

	return fmt.Sprintf("Unknown API: <b>%s %s", log.Method, log.Path)
}

var regPathTokenParser = regexp.MustCompile(`{([^}]+)}`)

func pathRegexp(path string) (*regexp.Regexp, []string) {
	tokens := []string{}
	src := regPathTokenParser.ReplaceAllStringFunc(path, func(s string) string {
		name := s[1 : len(s)-1]
		tokens = append(tokens, name)
		return "(?P<" + name + ">[^/]+)"
	})
	src = "^" + src + "$"
	re, _ := regexp.Compile(src)
	return re, tokens
}
*/
