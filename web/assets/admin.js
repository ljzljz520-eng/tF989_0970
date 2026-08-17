const elements = {
  rows: document.querySelector("#tag-rows"),
  empty: document.querySelector("#empty-state"),
  loadState: document.querySelector("#load-state"),
  menuStatus: document.querySelector("#menu-status"),
  tagCount: document.querySelector("#tag-count"),
  add: document.querySelector("#add-tag"),
  dialog: document.querySelector("#tag-dialog"),
  form: document.querySelector("#tag-form"),
  title: document.querySelector("#dialog-title"),
  id: document.querySelector("#tag-id"),
  name: document.querySelector("#tag-name"),
  color: document.querySelector("#tag-color"),
  sort: document.querySelector("#tag-sort"),
  scenes: document.querySelector("#tag-scenes"),
  error: document.querySelector("#form-error"),
  close: document.querySelector("#close-dialog"),
  cancel: document.querySelector("#cancel-dialog")
};

let tags = [];

function addTextCell(row, value, className = "") {
  const cell = document.createElement("td");
  const text = document.createElement("span");
  text.className = className;
  text.textContent = value;
  cell.append(text);
  row.append(cell);
  return cell;
}

function actionButton(label, className, handler) {
  const button = document.createElement("button");
  button.type = "button";
  button.className = className;
  button.textContent = label;
  button.addEventListener("click", handler);
  return button;
}

function renderRows() {
  elements.rows.replaceChildren();
  elements.empty.hidden = tags.length !== 0;
  elements.tagCount.textContent = `共 ${tags.length} 个标签`;

  for (const tag of tags) {
    const row = document.createElement("tr");
    addTextCell(row, String(tag.sortOrder));
    addTextCell(row, tag.name, "tag-name");

    const colorCell = document.createElement("td");
    const colorContent = document.createElement("span");
    colorContent.className = "color-cell";
    const swatch = document.createElement("span");
    swatch.className = "color-swatch";
    swatch.style.backgroundColor = tag.color;
    const colorText = document.createElement("span");
    colorText.textContent = tag.color;
    colorContent.append(swatch, colorText);
    colorCell.append(colorContent);
    row.append(colorCell);

    const sceneCell = document.createElement("td");
    const sceneList = document.createElement("span");
    sceneList.className = "scene-list";
    for (const scene of tag.applicableScenes) {
      const sceneLabel = document.createElement("span");
      sceneLabel.className = "scene-label";
      sceneLabel.textContent = scene;
      sceneList.append(sceneLabel);
    }
    sceneCell.append(sceneList);
    row.append(sceneCell);

    const actions = document.createElement("td");
    const actionGroup = document.createElement("span");
    actionGroup.className = "table-actions";
    actionGroup.append(
      actionButton("编辑", "link-button", () => openDialog(tag)),
      actionButton("删除", "danger", () => removeTag(tag))
    );
    actions.append(actionGroup);
    row.append(actions);
    elements.rows.append(row);
  }
}

function openDialog(tag) {
  elements.form.reset();
  elements.error.textContent = "";
  elements.id.value = tag?.id ?? "";
  elements.title.textContent = tag ? "编辑标签" : "新增标签";
  elements.name.value = tag?.name ?? "";
  elements.color.value = tag?.color ?? "#16855b";
  elements.sort.value = String(tag?.sortOrder ?? 10);
  elements.scenes.value = tag?.applicableScenes.join("，") ?? "";
  elements.dialog.showModal();
  elements.name.focus();
}

function closeDialog() {
  elements.dialog.close();
}

async function request(url, options = {}) {
  const response = await fetch(url, options);
  if (response.ok) {
    if (response.status === 204) {
      return null;
    }
    return response.json();
  }
  const payload = await response.json().catch(() => ({ error: "请求失败" }));
  const detail = payload.fields ? Object.values(payload.fields).join("；") : payload.error;
  throw new Error(detail || "请求失败");
}

async function refresh() {
  elements.loadState.textContent = "";
  try {
    const state = await request("/api/admin/state");
    tags = state.tags;
    elements.menuStatus.textContent = state.menus.some((menu) => menu.key === "customer-followup-labels")
      ? "插件菜单已安装"
      : "插件菜单未安装";
    renderRows();
  } catch (error) {
    elements.loadState.textContent = error.message;
  }
}

async function removeTag(tag) {
  if (!window.confirm(`确认删除“${tag.name}”？`)) {
    return;
  }
  try {
    await request(`/api/admin/tags/${encodeURIComponent(tag.id)}`, { method: "DELETE" });
    await refresh();
  } catch (error) {
    elements.loadState.textContent = error.message;
  }
}

elements.form.addEventListener("submit", async (event) => {
  event.preventDefault();
  elements.error.textContent = "";
  const id = elements.id.value;
  const input = {
    name: elements.name.value,
    color: elements.color.value,
    applicableScenes: elements.scenes.value.split(/[，,]/),
    sortOrder: Number(elements.sort.value)
  };
  try {
    await request(id ? `/api/admin/tags/${encodeURIComponent(id)}` : "/api/admin/tags", {
      method: id ? "PUT" : "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(input)
    });
    closeDialog();
    await refresh();
  } catch (error) {
    elements.error.textContent = error.message;
  }
});

elements.add.addEventListener("click", () => openDialog());
elements.close.addEventListener("click", closeDialog);
elements.cancel.addEventListener("click", closeDialog);
refresh();
