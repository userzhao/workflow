/*
 *  CAS 100credit
 */

bootbox.setLocale("zh_CN");

function loadgif(flag){
    var oLoading = $('#loading');
    if (oLoading){
        if(flag){
            oLoading.css('display', 'block');
        }else{
            oLoading.css('display', 'none');
        }
    }
}
