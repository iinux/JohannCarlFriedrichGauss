QUnit.test( "hello test", function( assert ) {
  assert.ok( 1 == "1", "Passed!" );
});

function format(string, values) {
    for (var key in values) {
        string = string.replace(new RegExp("\{" + key + "}","g"), values[key]);
    }
    return string;
}
 
test("basics", function() {
    var values = {
        name: "World"
    };
    equal( format("Hello, {name}", values), "Hello, World", "单个匹配" );
    equal( format("Hello, {name}, how is {name} today?", values),
        "Hello, World, how is World today?", "多个匹配" );
});