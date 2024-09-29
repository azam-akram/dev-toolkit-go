/*
 * Copyright 2021 ByteDance Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package avx

import (
    `github.com/bytedance/sonic/loader`
)


func Use() {
    loader.WrapGoC(_text_f64toa, _cfunc_f64toa, []loader.GoC{{"_f64toa", &S_f64toa, &F_f64toa}}, "avx", "avx/f64toa.c")
    loader.WrapGoC(_text_f32toa, _cfunc_f32toa, []loader.GoC{{"_f32toa", &S_f32toa, &F_f32toa}}, "avx", "avx/f32toa.c")
    loader.WrapGoC(_text_get_by_path, _cfunc_get_by_path, []loader.GoC{{"_get_by_path", &S_get_by_path, &F_get_by_path}}, "avx", "avx/get_by_path.c")
    loader.WrapGoC(_text_html_escape, _cfunc_html_escape, []loader.GoC{{"_html_escape", &S_html_escape, &F_html_escape}}, "avx", "avx/html_escape.c")
    loader.WrapGoC(_text_i64toa, _cfunc_i64toa, []loader.GoC{{"_i64toa", &S_i64toa, &F_i64toa}}, "avx", "avx/i64toa.c")
    loader.WrapGoC(_text_lspace, _cfunc_lspace, []loader.GoC{{"_lspace", &S_lspace, &F_lspace}}, "avx", "avx/lspace.c")
    loader.WrapGoC(_text_quote, _cfunc_quote, []loader.GoC{{"_quote", &S_quote, &F_quote}}, "avx", "avx/quote.c")
    loader.WrapGoC(_text_skip_array, _cfunc_skip_array, []loader.GoC{{"_skip_array", &S_skip_array, &F_skip_array}}, "avx", "avx/skip_array.c")
    loader.WrapGoC(_text_skip_number, _cfunc_skip_number, []loader.GoC{{"_skip_number", &S_skip_number, &F_skip_number}}, "avx", "avx/skip_number.c")
    loader.WrapGoC(_text_skip_one, _cfunc_skip_one, []loader.GoC{{"_skip_one", &S_skip_one, &F_skip_one}}, "avx", "avx/skip_one.c")
    loader.WrapGoC(_text_skip_object, _cfunc_skip_object, []loader.GoC{{"_skip_object", &S_skip_object, &F_skip_object}}, "avx", "avx/skip_object.c")
    loader.WrapGoC(_text_skip_one_fast, _cfunc_skip_one_fast, []loader.GoC{{"_skip_one_fast", &S_skip_one_fast, &F_skip_one_fast}}, "avx", "avx/skip_one_fast.c")
    loader.WrapGoC(_text_u64toa, _cfunc_u64toa, []loader.GoC{{"_u64toa", &S_u64toa, &F_u64toa}}, "avx", "avx/u64toa.c")
    loader.WrapGoC(_text_unquote, _cfunc_unquote, []loader.GoC{{"_unquote", &S_unquote, &F_unquote}}, "avx", "avx/unquote.c")
    loader.WrapGoC(_text_validate_one, _cfunc_validate_one, []loader.GoC{{"_validate_one", &S_validate_one, &F_validate_one}}, "avx", "avx/validate_one.c")
    loader.WrapGoC(_text_validate_utf8, _cfunc_validate_utf8, []loader.GoC{{"_validate_utf8", &S_validate_utf8, &F_validate_utf8}}, "avx", "avx/validate_utf8.c")
    loader.WrapGoC(_text_validate_utf8_fast, _cfunc_validate_utf8_fast, []loader.GoC{{"_validate_utf8_fast", &S_validate_utf8_fast, &F_validate_utf8_fast}}, "avx", "avx/validate_utf8_fast.c")
    loader.WrapGoC(_text_vnumber, _cfunc_vnumber, []loader.GoC{{"_vnumber", &S_vnumber, &F_vnumber}}, "avx", "avx/vnumber.c")
    loader.WrapGoC(_text_vsigned, _cfunc_vsigned, []loader.GoC{{"_vsigned", &S_vsigned, &F_vsigned}}, "avx", "avx/vsigned.c")
    loader.WrapGoC(_text_vunsigned, _cfunc_vunsigned, []loader.GoC{{"_vunsigned", &S_vunsigned, &F_vunsigned}}, "avx", "avx/vunsigned.c")
    loader.WrapGoC(_text_vstring, _cfunc_vstring, []loader.GoC{{"_vstring", &S_vstring, &F_vstring}}, "avx", "avx/vstring.c")
    loader.WrapGoC(_text_value, _cfunc_value, []loader.GoC{{"_value", &S_value, &F_value}}, "avx", "avx/value.c")
}
