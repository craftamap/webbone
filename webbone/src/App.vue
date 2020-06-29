<template>
    <div id="app">
        <form v-on:submit.prevent="submit" class="container-sm">
            <div class="mb-3">
                <label for="url" class="form-label">URL: </label>
                <input type="url" name="url" id="url" v-model="url" class="form-control"/>
            </div>
            <div class="mb-3">
                <label for="name" class="form-label">Name: </label>
                <input type="name" name="name" id="name" v-model="name" class="form-control"/>
            </div>
            <input type="submit" value="Create Thingy!" class="btn btn-primary"/>
            <div class="response" v-html="response">
            </div>
        </form>
    </div>
</template>

<script>


export default {
    name: 'App',
    data: function() {
        return {
            url: "",
            name: "",
            response: "",
        }
    },
    "methods": {
        "submit": function() {
            console.log("ola!")
            fetch("/r/", {
                "method": "POST",
                "body": JSON.stringify({
                    "url": this.url,
                    "name": this.name,
                })
            })
            .then(res => res.json())
            .then(data => {
                let url = "localhost:8080/r/"+ (data && data.name) || ""
                this.response = `<a href="${url}">${url}</a> is pointing to ${data.url}`
            })

        }
    }
}
</script>



<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
