const guide = document.querySelector("#tag-guide");
const state = document.querySelector("#guide-state");

function renderTag(tag) {
  const row = document.createElement("article");
  row.className = "guide-row";
  row.style.setProperty("--tag-color", tag.color);

  const order = document.createElement("span");
  order.className = "guide-order";
  order.textContent = String(tag.sortOrder).padStart(2, "0");

  const name = document.createElement("strong");
  name.className = "tag-name";
  name.textContent = tag.name;

  const scenes = document.createElement("span");
  scenes.className = "guide-scenes";
  for (const scene of tag.applicableScenes) {
    const label = document.createElement("span");
    label.className = "scene-label";
    label.textContent = scene;
    scenes.append(label);
  }

  const color = document.createElement("span");
  color.className = "guide-color";
  color.textContent = tag.color;
  row.append(order, name, scenes, color);
  return row;
}

fetch("/api/tags")
  .then((response) => {
    if (!response.ok) {
      throw new Error("标签读取失败");
    }
    return response.json();
  })
  .then((tags) => {
    state.hidden = tags.length > 0;
    state.textContent = tags.length === 0 ? "当前没有可用标签" : "";
    guide.replaceChildren(...tags.map(renderTag));
  })
  .catch((error) => {
    state.hidden = false;
    state.textContent = error.message;
  });
