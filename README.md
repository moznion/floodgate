floodgate
==

A tool to accumulate input from STDIN and flush that with according to rule.

Usage
--

### Options

```
Options:

  -h, --help              display help
  -i, --interval[=0]      intarval time to flush (second)
  -t, --threshold[=0]     throshold size of memory to flush (byte)
  -c, --concat[=\n]       character to concat for each line
      --stderr            flush to STDERR (default: STDOUT)
```

### Examples

#### Accumulate input and flush that every 30 seconds

```
$ tail -F /path/to/file | floodgate --interval=30
```

#### Accumulate input and flush that when buffer size is reached to 300 bytes

```
$ tail -F /path/to/file | floodgate --threshold=300
```

#### Accumulate input and flush that when buffer size is reached to 300 bytes or after waiting 30 seconds

```
$ tail -F /path/to/file | floodgate --interval=30 --threshold=300
```

Author
--

moznion (<moznion@gmail.com>)

License
--

MIT

