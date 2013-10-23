<?php 
    require_once("qiniu/rs.php");

    $domain = "aatest.qiniudn.com";
    $accessKey = "iN7NgwM31j4-BZacMjPrOQBs34UG1maYCAQmhdCV";
    $secretKey = "6QTOr2Jg1gcZEWDQXKOGZh5PziC2MCV5KsntT70j";
    
    $retStr = $_GET["upload_ret"];
    $errCode = $_GET["code"];
    $errMsg = urldecode($_GET["error"]);    
    
    if ($retStr)
        $decodedRet = base64_decode($retStr);
        $retObj = json_decode($decodedRet);
    
        $picKey = $retObj->{"key"};    
    
        Qiniu_SetKeys($accessKey, $secretKey);
        $baseUrl = Qiniu_RS_MakeBaseUrl($domain, $picKey);
        $getPolicy = new Qiniu_RS_GetPolicy();
        $privateUrl = $getPolicy->MakeRequest($baseUrl);            
?>

<html>
    <body>
        <p>UploadReult:</p>
        <?php 
            if ($retStr)
                echo "<p>$decodedRet</p>" . "<p>ImageDownloadUrl:<br>$privateUrl</p>" . "<p><img src=\"$privateUrl\"></p>";
            else
                echo "<p>error code: $errCode<br>error detail: $errMsg</p>";
        ?>        
        <p><a href="/upload.php">Back to Upload</a></p>
        <p><a href="/upload2.php">Back to UploadWithKey</a></p>
    </body>
</html>
