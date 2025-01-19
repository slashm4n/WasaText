<script setup>
import ErrorMsg from '../components/ErrorMsg.vue';
</script>

<script>
export default {
    emits: ['groupsUpdated', 'reloginNeeded'],
	props: ['session_token', 'my_groups'],
	data: function() {
		return {
            errormsg: '',
            selected_group: null,
            new_group_name: ''
		}
	},
	methods: {
        async doSetGroupName() {
            try {
                if (this.selected_group == null) {
                    this.errormsg = "Group name not set";
                    return;
                }

                if (this.new_group_name == '') {
                    this.errormsg = "New user group not set";
                    return;
                }
                
                const res = await this.$axios({
                    method: 'post',
                    url: '/groups/' + this.selected_group.group_id,
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

                this.$emit('groupsUpdated');
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
            try {
                const res = await this.$axios({
                    method: 'put',
                    url: '/group/photo',
                    headers: {
                        'Authorization' : 'Bearer ' + this.session_token
                    },
                    data: {
                        "group_id": this.selected_group.group_id,
                        "photo": img
                    }
                });

                if (res.status != 200) {
                    this.errormsg = "Problem while updating the group photo";
                }

                this.$emit('groupsUpdated');
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
                this.selected_group.photo = evn.target.result;
            };
        },

        async onPhotoGroupUploaderClick() {
            this.$refs.photoGroupUploader.value = ''
        },

        async doLeaveGroup() {
            try {
                if (this.selected_group == null) {
                    this.errormsg = "Group to leave not set";
                    return;
                }

                const res = await this.$axios({
                    method: 'delete',
                    url: '/groups/' + this.selected_group.group_id,
                    headers: {
                        'Authorization' : 'Bearer ' + this.session_token
                    },
                });

                if (res.status != 200) {
                    this.errormsg = "Problem while leaving the group";
                }

                this.$emit('groupsUpdated');
                this.errormsg = "";
            } catch (e) {
                if (e.response != null && e.response.data != "")
                    this.errormsg = "Error: " + e.response.data;
                else
                    this.errormsg = "Error: " + e;
            }
        }
    },
    watch: {
        session_token(newValue, oldValue) {
      		if (this.session_token == 0) {
                this.selected_group = null;
                this.new_group_name = '';
            }
		}
    },
    beforeMount: function () {
        window.addEventListener('beforeunload', (e) => {
            try {
                localStorage.setItem('selected_group',  JSON.stringify(this.selected_group));
                localStorage.setItem('new_group_name', JSON.stringify(this.new_group_name));
            } catch {
				this.$emit('reloginNeeded');
			}
        });

        try {
            if (localStorage.getItem('selected_group') != null) this.selected_group = JSON.parse(localStorage.getItem('selected_group'));
            if (localStorage.getItem('new_group_name') != null) this.new_group_name = JSON.parse(localStorage.getItem('new_group_name'));
        } catch {
        	this.$emit('reloginNeeded');
		}
    }
}
</script>

<template>
    <div v-if="session_token != 0">
        <div class="group-management-container">
            <div style="position:relative; top: 0.7em; float: left;">
	    		<span class="label-flat">For group</span>
                <select id="groupInput" style="position:relative; height: 1.3em; width: 9.5em;" v-model="selected_group" >
                    <option v-for="(g, index) in my_groups" :key="index" :value="g">{{ g.group_name }}</option>
                </select>                
                <input id="newGroupNameInput" v-model="new_group_name" placeholder="New group name" autocomplete="off" autocorrect="off" autocapitalize="off" spellcheck="false">
			    <button @click="doSetGroupName">Apply</button>
                <label for="photoGroupUploader" class="label-button">Set photo</label>
                <input id="photoGroupUploader" type="file" accept="image/*" hidden="true" ref="photoGroupUploader" @click="onPhotoGroupUploaderClick" @change="doSetGroupPhoto">
                <button @click="doLeaveGroup">Leave group</button>
            </div>
        </div>
        <ErrorMsg :errormsg="errormsg" @error-dismissed="this.errormsg = '';"></ErrorMsg>
    </div>
</template>
