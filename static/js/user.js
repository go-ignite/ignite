var User = function () {

    var signupHandler = function () {
        $('#signup-btn').on('click', function (e) {
            e.preventDefault();
            console.log("Clicked...")

            var form = $('#signup-form');
            var url = form.attr("action");
            $.post(url, form.serialize(), function (data) {
                if (data.success) {
                    console.log("Success...")
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