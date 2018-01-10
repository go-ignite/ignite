// client symbol
window.os = function() {  
    var ua = navigator.userAgent,  
    isWindowsPhone = /(?:Windows Phone)/.test(ua),  
    isSymbian = /(?:SymbianOS)/.test(ua) || isWindowsPhone,   
    isAndroid = /(?:Android)/.test(ua),   
    isFireFox = /(?:Firefox)/.test(ua),   
    isChrome = /(?:Chrome|CriOS)/.test(ua),  
    isTablet = /(?:iPad|PlayBook)/.test(ua) || (isAndroid && !/(?:Mobile)/.test(ua)) || (isFireFox && /(?:Tablet)/.test(ua)),  
    isPhone = /(?:iPhone)/.test(ua) && !isTablet,  
    isPc = !isPhone && !isAndroid && !isSymbian;  
    return {  
         isTablet: isTablet,  
         isPhone: isPhone,  
         isAndroid : isAndroid,  
         isPc : isPc  
    };
}();

 /*
  **********************************************************
  * OPAQUE NAVBAR SCRIPT
  **********************************************************
  */

  // Toggle tranparent navbar when the user scrolls the page

// $(window).scroll(function() {
//     if($(this).scrollTop() > 50)  /*height in pixels when the navbar becomes non opaque*/ 
//     {
//         $('.opaque-navbar').addClass('opaque');
//     } else {
//         $('.opaque-navbar').removeClass('opaque');
//     }
// });
