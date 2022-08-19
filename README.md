# Elasticsearch query string
[Elastic documentation](https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-query-string-query.html#query-string-syntax)

## Query string syntax
The query string “mini-language” is used by the Query string and by the q query string parameter in the search API.

The query string is parsed into a series of terms and operators. A term can be a single word — `quick or brown` — or a phrase, surrounded by double quotes — `"quick brown"` — which searches for all the words in the phrase, in the same order.

Operators allow you to customize the search — the available options are explained below.

## Field names
You can specify fields to search in the query syntax:

- where the status field contains active: `status:active`
- where the title field contains quick or brown: `title:(quick OR brown)`
- where the author field contains the exact phrase "john smith": `author:"John Smith"`
- where the first name field contains Alice (note how we need to escape the space with a backslash): `_exists_:title`

## Wildcards
Wildcard searches can be run on individual terms, using ? to replace a single character, and * to replace zero or more characters: `qu?ck bro*`

Be aware that wildcard queries can use an enormous amount of memory and perform very badly — just think how many terms need to be queried to match the query string "a* b* c*".

Allowing a wildcard at the beginning of a word (eg "*ing") is particularly heavy, because all terms in the index need to be examined, just in case they match. Leading wildcards can be disabled by setting allow_leading_wildcard to false.

Only parts of the analysis chain that operate at the character level are applied. So for instance, if the analyzer performs both lowercasing and stemming, only the lowercasing will be applied: it would be wrong to perform stemming on a word that is missing some of its letters.

By setting analyze_wildcard to true, queries that end with a * will be analyzed and a boolean query will be built out of the different tokens, by ensuring exact matches on the first N-1 tokens, and prefix match on the last token.

## Regular expressions
Regular expression patterns can be embedded in the query string by wrapping them in forward-slashes ("/"): `name:/joh?n(ath[oa]n)/`

The supported regular expression syntax is explained in Regular expression syntax.

The allow_leading_wildcard parameter does not have any control over regular expressions. A query string such as the following would force Elasticsearch to visit every term in the index: `/.*n/`
Use with caution!

## Fuzziness
You can run fuzzy queries using the ~ operator: `quikc~ brwn~ foks~`

For these queries, the query string is normalized. If present, only certain filters from the analyzer are applied. For a list of applicable filters, see Normalizers.

The query uses the Damerau-Levenshtein distance to find all terms with a maximum of two changes, where a change is the insertion, deletion or substitution of a single character, or transposition of two adjacent characters.

The default edit distance is 2, but an edit distance of 1 should be sufficient to catch 80% of all human misspellings. It can be specified as: `quikc~1`

Avoid mixing fuzziness with wildcards
Mixing fuzzy and wildcard operators is not supported. When mixed, one of the operators is not applied. For example, you can search for app~1 (fuzzy) or app* (wildcard), but searches for app*~1 do not apply the fuzzy operator (~1).

## Proximity searches
While a phrase query (eg "john smith") expects all of the terms in exactly the same order, a proximity query allows the specified words to be further apart or in a different order. In the same way that fuzzy queries can specify a maximum edit distance for characters in a word, a proximity search allows us to specify a maximum edit distance of words in a phrase: `"fox quick"~5`

The closer the text in a field is to the original order specified in the query string, the more relevant that document is considered to be. When compared to the above example query, the phrase "quick fox" would be considered more relevant than "quick brown fox".

## Ranges
Ranges can be specified for date, numeric or string fields. Inclusive ranges are specified with square brackets [min TO max] and exclusive ranges with curly brackets {min TO max}.

- All days in 2012: `date:[2012-01-01 TO 2012-12-31]`
- Numbers 1..5: `count:[1 TO 5]`
- Tags between alpha and omega, excluding alpha and omega: `tag:{alpha TO omega}`
- Numbers from 10 upwards: `count:[10 TO *]`
- Dates before 2012: `date:{* TO 2012-01-01}`
- Numbers from 1 up to but not including 5: `count:[1 TO 5}`

Ranges with one side unbounded can use the following syntax:

- `age:>10`
- `age:>=10`
- `age:<10`
- `age:<=10`

To combine an upper and lower bound with the simplified syntax, you would need to join two clauses with an AND operator:

- `age:(>=10 AND <20)`
- `age:(+>=10 +<20)`

The parsing of ranges in query strings can be complex and error prone. It is much more reliable to use an explicit range query.

## Boosting
Use the boost operator ^ to make one term more relevant than another. For instance, if we want to find all documents about foxes, but we are especially interested in quick foxes: `quick^2 fox`

The default boost value is 1, but can be any positive floating point number. Boosts between 0 and 1 reduce relevance.

Boosts can also be applied to phrases or to groups: `"john smith"^2   (foo bar)^4`

## Boolean operators
By default, all terms are optional, as long as one term matches. A search for foo bar baz will find any document that contains one or more of foo or bar or baz. We have already discussed the default_operator above which allows you to force all terms to be required, but there are also boolean operators which can be used in the query string itself to provide more control.

The preferred operators are + (this term must be present) and - (this term must not be present). All other terms are optional. For example, this query: `quick brown +fox -news` states that:

- fox must be present
- news must not be present
- quick and brown are optional — their presence increases the relevance

The familiar boolean operators AND, OR and NOT (also written &&, || and !) are also supported but beware that they do not honor the usual precedence rules, so parentheses should be used whenever multiple operators are used together. For instance the previous query could be rewritten as: `((quick AND fox) OR (brown AND fox) OR fox) AND NOT news`

This form now replicates the logic from the original query correctly, but the relevance scoring bears little resemblance to the original.
In contrast, the same query rewritten using the match query would look like this:

```
{
    "bool": {
        "must":     { "match": "fox"         },
        "should":   { "match": "quick brown" },
        "must_not": { "match": "news"        }
    }
}
```

## Grouping
Multiple terms or clauses can be grouped together with parentheses, to form sub-queries:

(quick OR brown) AND fox
Groups can be used to target a particular field, or to boost the result of a sub-query:

status:(active OR pending) title:(full text search)^2

## Reserved characters
If you need to use any of the characters which function as operators in your query itself (and not as operators), then you should escape them with a leading backslash. For instance, to search for `(1+1)=2` you would need to write your query as `\(1\+1\)\=2`. When using JSON for the request body, two preceding backslashes (\\) are required; the backslash is a reserved escaping character in JSON strings.

`GET /my-index-000001/_search`
```
{
  "query" : {
    "query_string" : {
      "query" : "kimchy\\!",
      "fields"  : ["user.id"]
    }
  }
}
```

 
The reserved characters are: `+ - = && || > < ! ( ) { } [ ] ^ " ~ * ? : \ /`

Failing to escape these special characters correctly could lead to a syntax error which prevents your query from running.

< and > can’t be escaped at all. The only way to prevent them from attempting to create a range query is to remove them from the query string entirely.


> *"Could" doesn't sound 100%, so I checked all and here are how we handle special characters*

- *100% failed characters are `( ) [ ] { }`*
- *`:` when first or last in the last term; this is too confusing, so I promoting : to 100% reserved*
- *`+` and `-` only when first character*
- *`!` only when first or last character*
- *`^ ~` will be converted to suffix*
- *`*` will try to convert to phrase*
- *`/` will try to convert to regex*
- *`* ?` are wildcard characters*
- *`\\` is escape character*
- *Our parser will try to convert to other term types on `( [ { " ^ ~`*
- *We escape first `+` and `-` and all `= & | > < ! ) } ] : /`*
	



## Whitespaces and empty queries
Whitespace is not considered an operator.

If the query string is empty or only contains whitespaces the query will yield an empty result set.