<script lang="ts" setup>
import {onMounted, reactive} from 'vue'
import {Project} from "../types";
import {OpenDirectoryDialog} from "../../wailsjs/go/main/App";
import {GetProjectInfo} from "../../wailsjs/go/main/Api";
import {v4 as uuidv4} from "uuid";
import ProjectBox from './ProjectBox.vue'
import {getProjectsList} from "../utils/storage";

const emits = defineEmits(["selected"])

const result : Project[] = [];
const data = reactive({
  folder: "",
  enableCreate: false,
  projects: result,
})

onMounted(() => {
  data.projects = getProjectsList()
})

function openDirectorySelectionDialog() {
  OpenDirectoryDialog().then(result => {
    data.folder = result
    data.enableCreate = true
  }).catch(error => {
    console.log("openDirectorySelectionDialog - error", error)
  })
}

function createProject() {
  if (data.folder === "") {
    console.log("Please select a project folder")
    return
  }

  GetProjectInfo(data.folder).then((projectInfo: any) => {
    if (projectInfo == null) {
      displayToast("alert", "The selected folder does not contains a valid Go Project")
      return
    }
    saveProject(projectInfo)
  }).catch((error: any) => {
    console.log("createProject - Error", error)
  })
}

function saveProject(projectInfo: any) {
  projectInfo.id = uuidv4()
  data.projects.push(projectInfo)
  localStorage.setItem('projects', JSON.stringify(data.projects));
  data.folder = ""
  data.enableCreate = false
  displayToast("success", "Go Project was created successfully")
}

function handleDelete({ id }: any) {
  data.projects = data.projects.filter(p => p.id !== id)
  localStorage.setItem('projects', JSON.stringify(data.projects));
}

function handleSelected(project: any) {
  emits("selected", project)
}

function canCreate() {
  return data.enableCreate
}

function displayToast(type: string, message: string) {
  let tClass = 'bg-info'
  if(type==='success'){
    tClass = 'bg-success'
  }
  if(type==='alert'){
    tClass = 'bg-danger'
  }
  (<any>$(document)).Toasts('create', {
    title: 'Error',
    body: message,
    class: tClass,
    fixed: false,
    autohide: true,
    delay: 6000
  })
}
</script>

<template>
  <div class="card">
    <div class="card-header ui-sortable-handle" style="cursor: move;">
      <h3 class="card-title">Registered Projects</h3>
      <div class="card-tools">
        <div class="btn btn-block btn-default btn-xs add-project" data-toggle="modal" data-target="#modal-create-project">
          <i class="fa fa-circle-plus"></i> Add Project
        </div>
      </div>
    </div>
    <div class="card-body">
      <div class="row" style="padding: 10px;">
        <project-box v-for="item in data.projects" :project=item @delete="handleDelete" @selected="(params: any) => handleSelected(params)" />
      </div>
    </div>
  </div>

  <div class="modal fade" id="modal-create-project" style="display: none;" aria-hidden="true">
    <div class="modal-dialog modal-lg">
      <div class="modal-content">
        <div class="modal-header">
          <h4 class="modal-title">Add Golang Project</h4>
          <button type="button" class="close" data-dismiss="modal" aria-label="Close">
            <span aria-hidden="true">??</span>
          </button>
        </div>
        <div class="modal-body">
          <div>
            <div class="folder-selection float-left">
              <span class="title">Go Project Directory: </span>
              <span v-if="data.folder === ''" class="path default">Please select a project folder</span>
              <span v-else class="path">{{ data.folder }}</span>
            </div>
            <button type="button" class="btn btn-default float-right"  v-on:click="openDirectorySelectionDialog">Choose Directory</button>
          </div>
        </div>
        <div class="modal-footer justify-content-between">
          <button type="button" class="btn btn-default" data-dismiss="modal">Cancel</button>
          <button type="button" class="btn btn-primary" v-bind:class="canCreate()? '':'disabled'" v-on:click="createProject" data-dismiss="modal">Create Project</button>
        </div>
      </div>

    </div>
  </div>
</template>

<style scoped>
.folder-selection {
  margin-top: 10px;
}
.folder-selection > .title {
  font-weight: bold;
}
.folder-selection > .path {
  font-size: 14px;
}
.folder-selection > .default {
  color: #aaaaaa;
}

.add-project {
  margin-right: 10px;
}
</style>
