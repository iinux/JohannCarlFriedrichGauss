<?php

function dd($var)
{
	var_dump($var);
	die(0);
}

$db_conn = mysqli_connect("localhost", "root", "", "test");

/* check connection */
if (mysqli_connect_errno()) {
    printf("Connect failed: %s\n", mysqli_connect_error());
    exit();
}

$provinces = array();

$query = 'SELECT prov FROM provinces ';

// Get query ---------------------------------------------------

$query_data = mysqli_query($db_conn, $query);

if ($query_data == FALSE)
    goto _adr1;

$rows = mysqli_num_rows($query_data);
if ($rows == 0)
{
    $provinces[] = 'None found';
    goto _adr1;
}

// Build array -------------------------------------------------

$provinces[] = 'All_';

while ($row = mysqli_fetch_array($query_data))
      $provinces[] = $row['prov'];

// End 

_adr1:
print json_encode($provinces);

mysqli_free_result($query_data);
mysqli_close($db_conn);
