<script setup>
import ErrorMsg from '../components/ErrorMsg.vue';
</script>

<script>
export default {
    emits: ['userAddedToGroup'],
	props: ['session_token'],
	data: function() {
		return {
            errormsg: '',
            group: null,
            user_name_to_add: '',
            group_name: '',
            new_group_name: ''
		}
	},
	methods: {
        async doAddToGroup() {
            try {
                if (this.user_name_to_add == '') {
                    this.errormsg = "User to add not set";
                    return;
                }

                if (this.group_name == '') {
                    this.errormsg = "Group name not set";
                    return;
                }
                
                const res = await this.$axios({
                    method: 'put',
                    url: '/groups',
                    headers: {
                        'Authorization' : 'Bearer ' + this.session_token
                    },
                    data: {
                        "user_name_to_add": this.user_name_to_add,
                        "group_name": this.group_name
                    }
                });

                if (res.status != 201) {
                    this.errormsg = "Unexpected response " + res.status;
                    return;
                }

                this.$emit('userAddedToGroup');
                this.user_name_to_add = ''
                this.group_name = '';
                this.errormsg = "";
            } catch (e) {
                if (e.response != null && e.response.data != "")
                    this.errormsg = "Error: " + e.response.data;
                else
                    this.errormsg = "Error: " + e;
            }
        },
        async doSetGroupName() {
            this.errormsg = "Not implemented";
            return;
            try {
                if (this.new_group_name == '') {
                    this.errormsg = "New user group not set";
                    return;
                }
                
                const res = await this.$axios({
                    method: 'put',
                    url: '/groups',
                    headers: {
                        'Authorization' : 'Bearer ' + this.session_token
                    },
                    data: {
                        "new_group_name": this.new_group_name
                    }
                });

                if (res.status != 200) {
                    this.errormsg = "Unexpected response " + res.status;
                    return;
                }

                // this.$emit('groupNameChanged');
                this.user_name_to_add = ''
                this.group_name = '';
                this.new_group_name = '';
                this.errormsg = "";
            } catch (e) {
                if (e.response != null && e.response.data != "")
                    this.errormsg = "Error: " + e.response.data;
                else
                    this.errormsg = "Error: " + e;
            }
        },

        // this code cannot be placed inside doSetMyPhoto because
        // await need to be on the top level of an async method
        async doUploadGroupPhoto(img)
        {
            this.errormsg = "Not implemented";
            return;
            try {
                const res = await this.$axios({
                    method: 'put',
                    url: '/user/photo',
                    headers: {
                        'Authorization' : 'Bearer ' + this.session_token
                    },
                    data: {
                        "photo": img
                    }
                });

                if (res.status != 200) {
                    this.errormsg = "Problem while updating the profile photo";
                }

                this.errormsg = "";
            } catch (e) {
                if (e.response != null && e.response.data != "")
                    this.errormsg = "Error: " + e.response.data;
                else
                    this.errormsg = "Error: " + e;
            }
        },
		
        async doSetGroupPhoto(e) {
            const img = e.target.files[0];
            const reader = new FileReader();
            reader.readAsDataURL(img);
            reader.onload = evn =>{
                this.doUploadGroupPhoto(evn.target.result);
                this.group.photo = evn.target.result;
            };
        }
    },
    watch: {
        session_token(newValue, oldValue) {
            if (newValue) {
                this.errormsg = '';
            }
        }
    }
}
</script>

<template>
    <div v-if="session_token != 0">
        <div class="group-management-container">
            <div style="position:relative; top: 0.7em; float: left;">
	    		<span class="label-flat">Group</span>
		    	<input v-model="group_name" autocomplete="off" placeholder="Group name"></input>
    			<input v-model="user_name_to_add" autocomplete="off" placeholder="New member"></input>
	    		<button @click="doAddToGroup">Apply</button>
                <input v-model="new_group_name" autocomplete="off" placeholder="New group name"></input>
			    <button @click="doSetGroupName">Apply</button>
                <label for="photoGroupUploader" class="label-button">Set photo</label>
                <input type="file" accept="image/*" hidden="true" id="photoGroupUploader" @change="doSetGroupPhoto">
            </div>
        </div>
        <ErrorMsg :errormsg="errormsg" @errorWindowClosed="this.errormsg = '';"></ErrorMsg>
    </div>
</template>
