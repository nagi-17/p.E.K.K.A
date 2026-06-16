import {create} from 'zustand'

export const useAuthStore=create(function(set){
    let token=localStorage.getItem('user_jwt');
    let userLogin=false;

    if (token!==null) {userLogin=true;}
    return { jwt: token, playerID: null, isAuthDone: userLogin,
        login: function(token, id){
            localStorage.setItem('user_jwt', token);
            set({jwt: token, playerID: id, isAuthDone: true});
        },
        logout: function() {
            localStorage.removeItem('user_jwt');
            set({jwt: null, playerID: null, isAuthDone: false});
        }
    };
});