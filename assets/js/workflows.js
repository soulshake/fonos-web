$(function(){
  var orig = $('.task')[0]
  var $appSelect = $(orig).find('.app')
  var repo_select_val = ""

  var buildTaskTMPL = _.template($("#buildTaskTMPL").html());
  var copyTaskTMPL = _.template($("#copyTaskTMPL").html());
  var runTaskTMPL = _.template($("#runTaskTMPL").html());
  var promoteTaskTMPL = _.template($("#promoteTaskTMPL").html());

  var $taskContainer = $("#tasks")

  var setRackId = function(el){
    var rackId = $(el).find(':selected').data('rack-id')
    $(el).closest(".task").find("input[name='rack_id']").val(rackId)
    $(el).closest(".task").find("input[name='rack_id_to']").val(rackId)
  }

  var taskListChanged = function(){
    $('#tasks .task').each(function(index){
      $(this).find('.task-number').html(index+1);
      $(this).find('select.app').each(function() {
        setRackId(this)
      })
    });
    if($("#tasks .task").length == 1) {
      $('.remove-task').hide();
    } else {
      $('.remove-task').show();
    }
  }

  //make list sortable
  $( "#tasks" ).sortable({
    stop: taskListChanged
  });
  $( "#tasks" ).disableSelection();

  var populateRepos = function(provider){
    var $repos = $("#repo")
    $repos.html("")
    switch(provider) {
      case 'github':
        $(githubRepos).each(function(index){
          $repos.append('<option value="' + this.id + '">' + this.full_name + "</option>")
        });
        break;
      case 'github-enterprise':
        $(githubEnterpriseRepos).each(function(index){
          $repos.append('<option value="' + this.id + '">' + this.full_name + "</option>")
        });
        break;
      case 'gitlab':
        $(gitlabRepos).each(function(index){
          $repos.append('<option value="' + this.id + '">' + this.full_name + "</option>")
        });
        break;
      default:
        throw new RuntimeException('Invalid provider');
    }
  }

  if(githubRepos.length > 0) {
    $("#integration").append('<option value="github">GitHub</option>')
  }

  if(githubEnterpriseRepos.length > 0) {
    $("#integration").append('<option value="github-enterprise">GitHub Enterprise Edition</option>')
  }

  if(gitlabRepos.length > 0) {
    $("#integration").append('<option value="gitlab">GitLab</option>')
  }

  populateRepos($("#integration").val())


  $("#integration").change(function(ev){
    ev.preventDefault()
    populateRepos($(this).val())
  });

  var appendTask = function(data) {
    if(data.kind === "build") {
      $taskContainer.append(buildTaskTMPL(data));
    } else if(data.kind === "copy") {
      $taskContainer.append(copyTaskTMPL(data));
    } else if(data.kind === "run") {
      $taskContainer.append(runTaskTMPL(data));
    } else if(data.kind === "promote") {
      $taskContainer.append(promoteTaskTMPL(data));
    } else {
      alert("Not Supported")
    }
  }

  //if this variable is set we are on the edit page and need to build the form
  if(tasks && tasks.length > 0) {
    populateRepos(integration.provider)
    $("select#integration").val(integration.provider)
    if(integration.provider == "gitlab") {
     $('#eventdiv').get(0).style.display ='none';
    }
    $("#repo").val(trigger.project_id)

    $("select#event").val(triggerEvent);
    if (triggerEvent == "pull_request") {
      $("#branchdiv").get(0).style.display ='none';
    }  else if (triggerEvent == "push") {
      $("#branchdiv").get(0).style.display ='block';
    }

    $(tasks).each(function(index){
      var task = this
      appendTask(this)
      var $newElement = $taskContainer.find(".task").last()
      _.each(task.params, function(val, key){
        var $input = $newElement.find("input[name='"+key +"']")
        if(key == "app_id"){ 
          var $input = $newElement.find("select[name='"+key +"']")
          //handle case where two racks have same app name
          $input.find("option[data-rack-id='"+ task.params.rack_id +"'][value='"+ val +"']").attr("selected","selected");
        } else if( key == "app_id_to") {
          var $input = $newElement.find("select[name='"+key +"']")
          $input.find("option[data-rack-id='"+ task.params.rack_id_to +"'][value='"+ val +"']").attr("selected","selected");
        } else {
          var $input = $newElement.find("input[name='"+key +"']")
          $input.val(val)
        }
      })
    })
  } else {
    console.log("here")
    console.log(racks)
    $taskContainer.append(buildTaskTMPL({racks: racks}));
  }

  taskListChanged()

  $(document).on('click', '#add-task', function(ev){
    ev.preventDefault();
    var type = $("#new_task_type").val()
    data = {
      racks: racks,
      kind: type
    }
    appendTask(data)
    taskListChanged();
  });

  $(document).on('click', '.remove-task', function(ev){
    ev.preventDefault();
    $(this).closest('.task').remove();
    taskListChanged();
  });

  $("#repo").on("select2:select", function (e) { 
    // $('.select-repo-warning').hide()
    // repo_select_val = $(e.currentTarget).val();
  });

  $(document).on("change", "select.app", function(e){
    setRackId(this)
  })

  $("#workflow-form").submit(function(ev){
    ////we are on the new page
    //if($("#repo").hasClass('select2') && repo_select_val === "") {
    //  window.scrollTo(0,0)
    //  $('.select-repo-warning').show()
    //  ev.preventDefault()
    //  return
    //} else {
    //  $('.select-repo-warning').hide()
    //}

    var tasks = [];
    $('.task').each(function(index){
      var task = {params:{}};
      $(this).find("input, select").each(function(i){
        var name = $(this).attr("name")
        if($(this).val() != ""){
          if(name == "kind") {
            task.kind = $(this).val()
          } else {
            task.params[name] = $(this).val()
          }
        }
      })

      $(this).find("input[type='checkbox']").each(function(i){
        task.parms[name] = $(this).is(":checked")
      })
      tasks.push(task)
    });
    $("#tasks_json").val(JSON.stringify(tasks))
  });
})
