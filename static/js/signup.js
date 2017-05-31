var Signup = function () {

    var signupHandler = function () {
        $('#signup-btn').on('click', function (e) {
            e.preventDefault();

            var form = $('#signup-form');
            var url = form.attr("action");
            $.post(url, form.serialize(), function (data) {
                if (data.result) {
                    window.location.href='/panel';
                } else {
                    //Signup failed
                    console.log('Signup failed!');
                    toastr.warning(data.message);
                    return false;
                }
            }, "json");
        });
    }

    return {
        //main function to initiate the module
        init: function () {
            signupHandler();
        }
    };
}();