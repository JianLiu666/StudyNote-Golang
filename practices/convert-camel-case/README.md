# Convert camel case to nsake case

read json file with camel case key field and write out converted key field name as snake_case.
so, what will you do to handle generic/arbitrary json format?

### Example 1:

given JSON
```JSON
{
    "version": "1.0",
    "rules": [
        {
            "resource": {
                "path": "/api/data/documents"
        },
        "allowOrigins": [ "http://this.example.com", "http://that.example.com" ],
        "allowMethods": [ "GET" ],
        "allowCredentials": true
        },
    ]
}
```

expected result
```JSON
{
    "version": "1.0",
    "rules": [
        {
            "resource": {
                "path": "/api/data/documents"
        },
        "allowOrigins": [ "http://this.example.com", "http://that.example.com" ],
        "allowMethods": [ "GET" ],
        "allowCredentials": true
        },
    ]
}
```

### Example 2:

given JSON
```JSON
{
    "stringField": "ShouldNotBeChanged",
    "numberField": 0,
    "stringArray": ["ShouldNotBeChanged", "ShouldNotBeChanged"],
    "numberArray": [0, 1],
    "nestedObj1": [
        {
            "fieldName11": 0,
            "fieldName12": "ShouldNotBeChanged"
        }
    ],
    "nestedObj2": {
        "fieldName21": 0,
        "fieldName22": "ShouldNotBeChanged",
        "nestedObj23": [
            {
                "fieldName231": 0,
                "fieldName232": "ShouldNotBeChanged"
            }
        ]
    },
    "nestedObj3": {
        "nestedObj31": {
            "nestedObj32": {
                "fieldName33": "ShouldNotBeChanged"
            }
        }
    }
}
```

expected result
```JSON
{
    "string_field": "ShouldNotBeChanged",
    "number_field": 0,
    "string_array": ["ShouldNotBeChanged", "ShouldNotBeChanged"],
    "number_array": [0, 1],
    "nested_obj1": [
        {
            "field_name11": 0,
            "field_name12": "ShouldNotBeChanged"
        }
    ],
    "nested_obj2": {
        "field_name21": 0,
        "field_name22": "ShouldNotBeChanged",
        "nested_obj23": [
            {
                "field_name231": 0,
                "field_name232": "ShouldNotBeChanged"
            }
        ]
    },
    "nested_obj3": {
        "nested_obj31": {
            "nested_obj32": {
                "field_name33": "ShouldNotBeChanged"
            }
        }
    },
}
```

### Constraints:

- using stand library only. don’t use any 3rd party library/framework
- allowed language: golang