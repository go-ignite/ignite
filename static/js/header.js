
var Header = function () {
    return {
        //main function to initiate the module
        init: function () {
           	$('.header-li').on('click', function (e) {
				$(this).addClass("active").siblings().removeClass("active");
			});
        }
    };
}();