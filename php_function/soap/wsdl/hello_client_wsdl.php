<?php
try {
    $client = new SoapClient('hello.wsdl');
    
    $result =  $client->__soapCall('greet', [
        ['name' => 'Suhua']
    ]);
    
    printf("Result = %s", $result->greetReturn);
} catch (Exception $e) {
    printf("Message = %s", $e->__toString());
}