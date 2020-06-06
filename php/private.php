<?php

class LockedGate
{
    private function open()
    {
        return 'how did you get in here?!!' ."\n";
    }
    private function open2($a)
    {
        return "how did you get in here [$a] ?!! \n";
    }
}

$object = new LockedGate();
$reflector = new ReflectionObject($object);

$method = $reflector->getMethod('open');
$method->setAccessible(true);
echo $method->invoke($object);

$method2 = $reflector->getMethod('open2');
$method2->setAccessible(true);
echo $method2->invoke($object, "ha");
echo $method2->invokeArgs($object, ["ha"]);
// =============================
class Foo {
    private $bar = "Foo::Bar";
    private function add_ab($a, $b) {
        return $a + $b;
    }
}

$foo = new Foo;

// Single variable example
$getFooBarCallback = function() {
    return $this->bar;
};

$getFooBar = $getFooBarCallback->bindTo($foo, 'Foo');

echo $getFooBar()."\n"; // Prints Foo::Bar

// Function call with parameters example
$getFooAddABCallback = function() {
    // As of PHP 5.6 we can use $this->fn(...func_get_args()) instead of call_user_func_array
    return call_user_func_array(array($this, 'add_ab'), func_get_args());
};

$getFooAddAB = $getFooAddABCallback->bindTo($foo, 'Foo');

echo $getFooAddAB(33, 6)."\n"; // Prints 39
// =============================
$foo = new Foo;

// Single variable example
$getFooBar = function() {
    return $this->bar;
};

echo $getFooBar->call($foo)."\n"; // Prints Foo::Bar

// Function call with parameters example
$getFooAddAB = function() {
    return $this->add_ab(...func_get_args());
};

echo $getFooAddAB->call($foo, 33, 6)."\n"; // Prints 39
// =============
// refer https://stackoverflow.com/questions/2738663/call-private-methods-and-private-properties-from-outside-a-class-in-php/20333436
