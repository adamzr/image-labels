/*
 * Copyright 2018-2020 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package labels

import (
	"fmt"
	"os"
	"strings"

	"github.com/buildpacks/libcnb"
	"github.com/mattn/go-shellwords"
	"github.com/paketo-buildpacks/libpak/bard"
)

type Build struct {
	Logger bard.Logger
}

func (b Build) Build(context libcnb.BuildContext) (libcnb.BuildResult, error) {
	b.Logger.Title(context.Buildpack)
	b.Logger.Body(bard.FormatUserConfig("BP_IMAGE_LABELS", "arbitrary image labels", "<none>"))
	for k, v := range Labels {
		b.Logger.Body(bard.FormatUserConfig(k, fmt.Sprintf("the %s image label", v), "<none>"))
	}

	result := libcnb.BuildResult{}

	for k, v := range Labels {
		if s, ok := os.LookupEnv(k); ok {
			result.Labels = append(result.Labels, libcnb.Label{Key: v, Value: s})
		}
	}

	if s, ok := os.LookupEnv("BP_IMAGE_LABELS"); ok {
		words, err := shellwords.Parse(s)
		if err != nil {
			return libcnb.BuildResult{}, fmt.Errorf("unable to parse %s\n%w", s, err)
		}

		for _, word := range words {
			parts := strings.Split(word, "=")
			result.Labels = append(result.Labels, libcnb.Label{Key: parts[0], Value: parts[1]})
		}
	}

	return result, nil
}
