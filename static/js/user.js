var User = function () {

    var signupHandler = function () {
        $(".signup").click(function(event){
            event.preventDefault();
            $('#hero').fadeOut(1000);
            $('#hero').css("display", "none");

            $('#login').css("display", "none")
            $('#signup').fadeIn(500);
        });

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
        $('.login').on('click', function (e) {
            e.preventDefault();
            $('#hero').fadeOut(1000);
            $('#hero').css("display", "none");

            $('#signup').css("display", "none")
            $('#login').fadeIn(500);
        });

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

    return {
        //main function to initiate the module
        init: function () {
            signupHandler();
            loginHandler();
        }
    };
}();