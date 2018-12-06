new Vue({
    el: '#app',

    data: {
        ws: null, // Our websocket
        newMsg: '', // Holds new messages to be sent to the server
        chatContent: '', // A running list of chat messages displayed on the screen
        username: null,
        email: null,
        joined: false // True if email and username have been filled in
    },
    created: function () {
        var self = this;
        this.ws = new WebSocket('ws://' + window.location.host + '/ws');
        this.ws.addEventListener('message', function (e) {
            var msg = JSON.parse(e.data)
            self.chatContent += '<div class="chip">' +
                '<img src="' + self.gravatarURL(msg.user.email) + '">' // Avatar
                +
                msg.user.username +
                '</div>' +
                emojione.toImage(msg.message) + '<br/>'; // Parse emojis

            var element = document.getElementById('chat-messages');
            element.scrollTop = element.scrollHeight; // Auto scroll to the bottom
        });
        this.ws.addEventListener('onclose', function(e) {
            
        });
    },
    methods: {
        send: function () {
            if (this.newMsg != '') {
                this.ws.send(
                    JSON.stringify({
                        user: {
                            email: this.email,
                            username: this.username
                        },
                        type: 0,
                        message: $('<p>').html(this.newMsg).text() // Strip out html
                    }));
                this.newMsg = ''; // Reset newMsg
            }
        },
        join: function () {
            if (!this.email) {
                Materialize.toast('You must enter an email', 2000);
                return
            }
            if (!this.username) {
                Materialize.toast('You must choose a username', 2000);
                return
            }
            this.email = $('<p>').html(this.email).text();
            this.username = $('<p>').html(this.username).text();
            this.joined = true;
            this.ws.send(
                JSON.stringify({
                    user: {
                        email: this.email,
                        username: this.username
                    },
                    type: 1
                }));
        },
        gravatarURL: function (email) {
            return 'http://www.gravatar.com/avatar/' + CryptoJS.MD5(email);
        }
    }
});