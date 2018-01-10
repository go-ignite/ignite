var User = function () {
    var currentId = 'hero';

    function animateToggle(targetId) {
        $('#' + currentId).fadeOut(500);
        $('#' + currentId).css("display", "none");
        $('.header-' + targetId).addClass("active").siblings().removeClass("active");
        $('#' + targetId).fadeIn(1000);
        currentId = targetId; 
    }
    var signupHandler = function () {

        $('#signup-btn').on('click', function (e) {
            e.preventDefault();

            var form = $('#signup-form');
            var url = form.attr("action");
            $.post(url, form.serialize(), function (data) {
                if (data.success) {
                    window.location.href = '/panel/index';
                } else {
                    //Signup failed
                    toastr.warning(data.message);
                    return false;
                }
            }, "json");
        });
    };

    var loginHandler = function () {

        $('#login-btn').on('click', function (e) {
            e.preventDefault();

            var form = $('#login-form');
            var url = form.attr("action");
            $.post(url, form.serialize(), function (data) {
                if (data.success) {
                    console.log("login successfully!");
                    console.log(data);

                    window.location.href = '/panel/index';
                } else {
                    //Login failed
                    toastr.warning(data.message);
                    return false;
                }
            }, "json");
        });
    };

    var commonHandler = function () {
        $('.header-li').on('click', function (e) {
            $(this).addClass("active").siblings().removeClass("active");
        });
        $(".hero-move").click(function(){
            animateToggle('hero');
        });
        $(".signup-move").click(function(){
            animateToggle('signup');
        });
        $('.login-move').on('click', function() {
            animateToggle('login');
        });
    };

    return {
        //main function to initiate the module
        init: function () {
            signupHandler();
            loginHandler();
            commonHandler();
        }
    };
}();
