<?php 
    require_once("qiniu/rs.php");
    
    $bucket = "a";
    $accessKey = "iN7NgwM31j4-BZacMjPrOQBs34UG1maYCAQmhdCV";
    $secretKey = "6QTOr2Jg1gcZEWDQXKOGZh5PziC2MCV5KsntT70j";
    
    Qiniu_SetKeys($accessKey, $secretKey);
    $putPolicy = new Qiniu_RS_PutPolicy($bucket);
    $putPolicy->ReturnUrl = "http://localhost/uploaded.php";
    $upToken = $putPolicy->Token(null);
    #echo $upToken;
?>
<html>
    <body>
        <form method="post" action="http://up.qiniu.com/" enctype="multipart/form-data">
            <input name="token" type="hidden" value="<?php 
                print $upToken
            ?>">
            Album belonged to:
            <input type="text" name="x:album" value="albumId"><br>
            Image to upload:
            <input type="file" name="file"><br>
            <button type="submit">Upload</button>
        </form>
    </body>
</html>
