# Wrench
Library for more easily using [BBolt](https://github.com/etcd-io/bbolt)
Based almost entirely on internal methods from my [BoltInspector](https://github.com/89yoyos/BoltInspector) tool

## About
Wrench is a tool I am working on to simplify implementing Bolt databases. The idea is to make reading the database as easy as simply using Insert and Delete commands (among others). Obviously, majorly simplifying the library also limits it in some ways, however I feel the tradeoffs are worthwhile for developers (like myself) who aren't stressing over perfect implementation and simply want a quick and easy way to store information.

## Usage
### Adding Wrench to Your Project

You can grab this library with Go Get
```
$ go get https://github.com/89yoyos/Wrench/...
```

Then import the library like any other
```go
import ("https://github.com/89yoyos/Wrench")
```

### Setup
To use Wrench, you need to declare a Wrench object.
```go
w := Wrench.Wrench{}
```
Next, just use your Wrench to open the Bolt database you want to use.
```go
err := w.OpenDB("File/Path.db")
```
The `OpenDB` command returns a normal Error object if it fails to open the file, so you can implement whatever error handling you're most comfortable with.

### Adding a Bucket
An empty Bolt database cannot hold key value pairs, they need to be inside of a bucket. To create a bucket with Wrench, use the `CreateBucket` command.
```go
w.CreateBucket("key")
```
If you've used Bolt before, you may have noticed that Wrench treats all keys as strings. This is intentional, though not in line with pure Bolt, which allows any array of bytes to be used as a key. I found this limitation immensely simplifies working with the database, though it may cause issues with databases with non-string keys.

### Navigating Your Database
In Bolt, accessing buckets could be tedious, and accessing sub-buckets would only make it worse. Wrench attempts to alleviate this by treating the database like a filesystem. Your Wrench object knows its current location within the database, and anything you do is relative to position. To see where your Wrench is currently positioned, you can use the `CurrentBucket` and `CurrentBucketString` functions.

`CurrentBucket` returns an array of strings, starting at the root of the directory and working through all the buckets and sub-buckets to get to your current position.

`CurrentBucketString` works in a similar way, however it returns a single string representing that path.

To change where your Wrench is positioned, you can use the `GoTo` function.
```go
w.GoTo(["~","Path","as","an","Array"])
```
Or, if you don't want to construct an array, you can use the built-in method to do that
```go
w.GoToString("~/Path/as/a/String")
```
This function uses the `StringToPath` function to construct an array of string, and passes that to the `GoTo` function. The benefit being that the `StringToPath` function does some nifty things, like allowing relative paths to be used.

Additionally, note the `~` in both examples. A `~` denotes the root of the database for absolute paths. If you use the `GoToString` method or the `StringToPath` method, it will create an array that includes the `~`, adding it to the path even if you forget to include it.

### Inserting Values
Bolt does not support inserting values into the root of the database. So to insert a value, you must first create a bucket, then navigate your wrench into that bucket. With Wrench, this is very simple to do.
```go
w := Wrench.Wrench{}
err := w.OpenDB("File/Path.db")
// Error Handling Here
w.CreateBucket("Example")
w.GoToString("~/Example")
```
In 4 lines of code, you created a Wrench, opened the database, created a bucket and are now ready to write to that bucket.

Inserting a value is just as easy.
```go
w.Insert("key",[value as bytes])
```
Like buckets, values' keys are limited to strings when using Wrench. Their values are arrays of bytes. To convert a value to bytes, Go has built in methods.
```go
w.Insert("key",[]byte("value"))
```
This method also works with non-string values.

### Reading Values
Because of how Bolt works, the default method to read values from a bucket is to get all of them at once and sort through them. This can be done with the `GetAll` function.
```go
w.GetAll()
```

This function, along with all others that read from the database, return a DBVal object. DBVals are just things in the database. They can be either Values or Buckets. They include the `path` to the value, the `key` for that value, and the `value` itself. If a DBVal represents a bucket in the database, the `value` will be `nil`.

DBVals can be used to manipulate entries in the database, however that functionality is still under development.

If you don't want all of the values, the `GetOne` function allows you to get just the one value you're looking for.
```go
w.Get("key")
```

There are also several other `Get` functions.
```go
w.GetValues() // Returns just the non-bucket values in the current bucket
w.GetBuckets() // Returns just the buckets in the current bucket
w.GetBoth() // Returns 2 arrays, one of values and one of buckets
```

There are also function to count the values in the current bucket.
```go
w.Count() // Counts everything in the bucket
w.CountValues() // Counts the non-bucket values
w.CountBuckets() // Counts the buckets
w.CountBoth() // Returns the counts of the buckets and of the non-buckets as separate numbers
```

### Deleting Values
There are 2 commands that can be used to delete database entries.
```go
w.Delete("key")
```
The `Delete` command can delete both values and buckets, and it works exactly how you would expect.

```go
w.Empty()
```
The `Empty` command is a bit less intuitive, however it can be very useful in some circumstances. This command deletes all values and buckets in the Wrench's current location. For example, let's say your database looks like this:
```
~ [Root]
- TestBucket
- - TestVal1
- - TestVal2
```
If you point your Wrench at `~/TestBucket` and run `w.Empty()`, it will delete `~/TestBucket/TestVal1` and `~/TestBucket/TestVal2`, but `~/TestBucket` will remain in place, now empty of any values.

## Conclusion
Those are the basics of using Wrench. It's simple, but powerful. Additional functionality will likely be added in the future as I use this library in my own projects.

Among planned updates is the expansion of the `DBVal` class. I'd like to make it a more useful tool for manipulating data rather than simply reading it. The `DBVal.AsWrench()` function is a working, but unsatisfying version of this for now, which returns a new Wrench object that can then be used to manipulate data. Future versions will likely add functions to `DBVal` to update and delete the referenced entries.
