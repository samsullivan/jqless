# jqless

`jqless` is a combination between `jq` and `less`, enabling users to filter and extract data from JSON in real-time -- useful when first learning the syntax or for power users trying to extract multiple pieces of data from a single JSON blob.

![demo](https://github.com/samsullivan/jqless/blob/main/assets/demo.gif?raw=true)

#### Usage

To use, start `jqless` in your favorite terminal by either piping JSON data to the process or including a file path as the first argument.

```
cat /path/to/file.json | jqless
jqless /path/to/file.json
```

Once loaded, type your `jq` query as expected and see the results filter. To extract results to your clipboard, use `ctrl+x` as shown in help text.

More options to come in future versions!

#### Acknowledgements

It is written in Golang using the [Bubble Tea framework](https://github.com/charmbracelet/bubbletea) and [`gojq`](https://github.com/itchyny/gojq). Inspiration h/t to [`jq-live`](https://github.com/TheDahv/jq-live).