//    Copyright 2021. Go-Ceres
//    Author https://github.com/go-ceres/go-ceres
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package uri

import "regexp"

var SwaggerUrlRegexp = regexp.MustCompile("{[.;?]?([^{}*]+)\\*?}")

// SwaggerUriToGinUri 转gin路由
func SwaggerUriToGinUri(uri string) string {
	return SwaggerUrlRegexp.ReplaceAllString(uri, ":$1")
}

// SwaggerUriToFiberUri 转fiber路由
func SwaggerUriToFiberUri(uri string) string {
	return SwaggerUrlRegexp.ReplaceAllString(uri, ":$1")
}

// SwaggerUriToEchoUri 转echo路由
func SwaggerUriToEchoUri(uri string) string {
	return SwaggerUrlRegexp.ReplaceAllString(uri, ":$1")
}

func OrderedParamsFromUri(uri string) []string {
	matches := SwaggerUrlRegexp.FindAllStringSubmatch(uri, -1)
	result := make([]string, len(matches))
	for i, m := range matches {
		result[i] = m[1]
	}
	return result
}
