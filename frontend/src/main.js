const ANY = 255;

let ROLE_LABELS = {
  0: "铁甲战士",
  1: "静默猎手",
  2: "摄政",
  3: "缚灵师",
  4: "故障机器人",
  5: "无色",
  255: "全部角色",
};

let TYPE_LABELS = {
  0: "攻击",
  1: "技能",
  2: "能力",
  3: "任务",
  4: "状态",
  5: "诅咒",
  255: "全部类型",
};

let COST_LABELS = {
  0: "0费",
  1: "1费",
  2: "2费",
  3: "3费",
  4: "X费",
  255: "全部费用",
};

let RARITY_LABELS = {
  0: "基础",
  1: "普通",
  2: "罕见",
  3: "稀有",
  4: "远古",
  255: "全部稀有度",
};

const UPGRADE_LABELS = {
  true: "基础卡(未升级)",
  false: "升级卡",
};

const DECK_UPGRADE_FILTER_LABELS = {
  all: "全部形态",
  base: "基础卡(未升级)",
  upgraded: "升级卡",
};

let TYPE_LABELS_DETAIL = { ...TYPE_LABELS };

const state = {
  selectedDeckId: null,
  selectedDeckRole: ANY,
  selectedDeckName: "",
  completionNotified: false,
  completionToastTimer: null,
};

function backend() {
  return window?.go?.app?.App ?? window?.go?.main?.App ?? null;
}

async function invoke(method, ...args) {
  const api = backend();
  if (!api || typeof api[method] !== "function") {
    throw new Error(`Wails 绑定未就绪，方法不可用: ${method}`);
  }
  return api[method](...args);
}

async function runAction(action, fn) {
  try {
    await fn();
  } catch (err) {
    console.error(`[${action}]`, err);
    alert(`${action}失败: ${err?.message ?? String(err)}`);
  }
}

function setOptions(select, mapping, selected = String(ANY)) {
  select.innerHTML = "";
  Object.entries(mapping).forEach(([value, label]) => {
    const option = document.createElement("option");
    option.value = value;
    option.textContent = label;
    if (String(value) === String(selected)) {
      option.selected = true;
    }
    select.appendChild(option);
  });
}

function labelsFromItems(items, fallbackMap) {
  const next = { ...fallbackMap };
  if (!Array.isArray(items)) {
    return next;
  }
  items.forEach((item) => {
    if (item && typeof item.value !== "undefined" && typeof item.label === "string") {
      next[String(item.value)] = item.label;
    }
  });
  return next;
}

async function loadDisplayLabels() {
  try {
    const labels = await invoke("GetDisplayLabels");
    if (!labels || typeof labels !== "object") {
      return;
    }
    ROLE_LABELS = labelsFromItems(labels.roles, ROLE_LABELS);
    TYPE_LABELS = labelsFromItems(labels.types, TYPE_LABELS);
    TYPE_LABELS_DETAIL = { ...TYPE_LABELS };
    COST_LABELS = labelsFromItems(labels.costs, COST_LABELS);
    RARITY_LABELS = labelsFromItems(labels.rarities, RARITY_LABELS);
  } catch (err) {
    console.warn("load display labels failed, fallback to defaults", err);
  }
}

function toImagePath(image) {
  if (!image) return "";
  if (image.startsWith("/")) return image;
  return `/${image}`;
}

function showEmpty(container, text) {
  container.innerHTML = `<div class="empty">${text}</div>`;
}

function resetDeckTypeFilters() {
  document.getElementById("deckTypeKeyword").value = "";
  document.getElementById("deckTypeFilter").value = String(ANY);
  document.getElementById("deckTypeUpgradeFilter").value = "all";
}

function resetCardFilters() {
  document.getElementById("keyword").value = "";
  document.getElementById("cardRoleFilter").value = String(ANY);
  document.getElementById("typeFilter").value = String(ANY);
  document.getElementById("costFilter").value = String(ANY);
  document.getElementById("rarityFilter").value = String(ANY);
  document.getElementById("upgradeFilter").value = "true";
}

function typeLabel(typeValue) {
  return TYPE_LABELS_DETAIL[typeValue] ?? String(typeValue);
}

function roleLabel(roleValue) {
  return ROLE_LABELS[roleValue] ?? String(roleValue);
}

function rarityLabel(rarityValue) {
  return RARITY_LABELS[rarityValue] ?? String(rarityValue);
}

function costLabel(costValue) {
  return COST_LABELS[costValue] ?? String(costValue);
}

function showCompletionToast() {
  const toast = document.getElementById("completionToast");
  if (!toast) return;
  toast.classList.remove("hidden");
  toast.classList.remove("show");
  if (state.completionToastTimer) {
    clearTimeout(state.completionToastTimer);
    state.completionToastTimer = null;
  }
  void toast.offsetWidth;
  toast.classList.add("show");
  state.completionToastTimer = setTimeout(() => {
    toast.classList.remove("show");
    toast.classList.add("hidden");
  }, 3000);
}

function updateCompletionLabel(completed) {
  const label = document.getElementById("deckCompletionLabel");
  if (!label) return;
  if (completed) {
    label.classList.remove("hidden");
  } else {
    label.classList.add("hidden");
  }
}

async function refreshCompletionStatus() {
  if (!state.selectedDeckId) {
    state.completionNotified = false;
    updateCompletionLabel(false);
    return;
  }
  const status = await invoke("GetDeckCompletionStatus");
  const completed = Boolean(status?.completed);
  updateCompletionLabel(completed);
  if (completed && !state.completionNotified) {
    showCompletionToast();
    state.completionNotified = true;
  }
  if (!completed) {
    state.completionNotified = false;
  }
}

async function refreshDeckTypeStats() {
  const container = document.getElementById("deckTypeStats");
  if (!state.selectedDeckId) {
    showEmpty(container, "未选择牌组");
    return;
  }

  const keyword = document.getElementById("deckTypeKeyword").value.trim();
  const type = Number(document.getElementById("deckTypeFilter").value);
  const upgradeFilter = document.getElementById("deckTypeUpgradeFilter").value;

  const statsRaw = await invoke("GetSelectedDeckCardStats", type, keyword, upgradeFilter);
  const stats = Array.isArray(statsRaw) ? statsRaw : [];
  if (!stats.length) {
    showEmpty(container, "该筛选条件下没有卡牌");
    return;
  }

  const grouped = {};
  stats.forEach((entry) => {
    const key = String(entry.type);
    if (!grouped[key]) {
      grouped[key] = {
        type: entry.type,
        total: 0,
        cards: [],
      };
    }
    grouped[key].total += Number(entry.count || 0);
    grouped[key].cards.push(entry);
  });

  const orderedTypes = Object.keys(grouped)
    .map((value) => Number(value))
    .sort((a, b) => a - b);

  container.innerHTML = "";
  orderedTypes.forEach((typeValue) => {
    const group = grouped[String(typeValue)];
    const block = document.createElement("div");
    block.className = "type-group";
    block.innerHTML = `
      <div class="type-group-head">
        <strong>${typeLabel(typeValue)}</strong>
        <span class="badge">总数 ${group.total}</span>
      </div>
      <div class="list type-group-body"></div>
    `;

    const body = block.querySelector(".type-group-body");
    group.cards.forEach((card) => {
      const row = document.createElement("div");
      row.className = "item item-compact";
      row.innerHTML = `
        <div class="item-meta">
          <img src="${toImagePath(card.image)}" alt="${card.id}" data-card-id="${card.id}" />
          <div>
            <strong>${card.name || card.id}</strong>
            <div class="item-subtext">${card.id}</div>
            <div>
              <span class="badge">${card.is_base ? "基础卡" : "升级卡"}</span>
              <span class="badge need">x${card.count}</span>
            </div>
          </div>
        </div>
      `;
      const image = row.querySelector("img[data-card-id]");
      if (image) {
        image.addEventListener("click", () => {
          openCardDetail(card.id);
        });
      }
      body.appendChild(row);
    });

    container.appendChild(block);
  });
}

async function openCardDetail(cardID) {
  await runAction("查看卡牌详情", async () => {
    const card = await invoke("GetCardDetail", cardID);
    document.getElementById("detailImage").src = toImagePath(card.image);
    document.getElementById("detailImage").alt = card.id;
    document.getElementById("detailName").textContent = card.name || card.id || "-";
    document.getElementById("detailID").textContent = `ID: ${card.id ?? "-"}`;
    document.getElementById("detailAlias").textContent = `别名: ${card.alias ?? "-"}`;
    document.getElementById("detailType").textContent = `类型: ${typeLabel(card.type)}`;
    document.getElementById("detailRole").textContent = `角色: ${roleLabel(card.role)}`;
    document.getElementById("detailCost").textContent = `费用: ${costLabel(card.cost)}`;
    document.getElementById("detailRarity").textContent = `稀有度: ${rarityLabel(card.rarity)}`;
    document.getElementById("detailUpgrade").textContent = `形态: ${card.is_base ? "基础卡(未升级)" : "升级卡"}`;

    const modal = document.getElementById("cardDetailModal");
    modal.classList.remove("hidden");
    modal.setAttribute("aria-hidden", "false");
  });
}

function closeCardDetail() {
  const modal = document.getElementById("cardDetailModal");
  modal.classList.add("hidden");
  modal.setAttribute("aria-hidden", "true");
}

async function refreshDeckList() {
  const role = Number(document.getElementById("deckRoleFilter").value);
  const deckListRaw = await invoke("GetDeckList", role);
  const deckList = Array.isArray(deckListRaw) ? deckListRaw : [];
  const container = document.getElementById("deckList");
  if (!deckList.length) {
    showEmpty(container, "无可用牌组");
    return;
  }

  container.innerHTML = "";
  deckList.forEach((deck) => {
    const row = document.createElement("div");
    row.className = "item";
    const deckTitle = deck.name?.trim() ? deck.name : "未命名牌组";
    row.innerHTML = `
      <div class="item-meta">
        <div>
          <strong>${deckTitle}</strong>
          <div class="badge">${ROLE_LABELS[deck.role] ?? deck.role}</div>
        </div>
      </div>
    `;
    const btn = document.createElement("button");
    const isSelected = state.selectedDeckId === deck.id;
    if (isSelected) {
      row.classList.add("selected");
    }
    btn.textContent = isSelected ? "取消选择" : "选择";
    btn.addEventListener("click", async () => {
      await runAction("选择牌组", async () => {
        await invoke("ChooseDeck", deck.id);
        if (isSelected) {
          state.selectedDeckId = null;
          state.selectedDeckRole = ANY;
          state.selectedDeckName = "";
          state.completionNotified = false;
          document.getElementById("cardRoleFilter").value = String(ANY);
          document.getElementById("selectedDeck").textContent = "当前牌组: 未选择";
        } else {
          state.selectedDeckId = deck.id;
          state.selectedDeckRole = deck.role;
          state.selectedDeckName = deckTitle;
          state.completionNotified = false;
          document.getElementById("cardRoleFilter").value = String(deck.role);
          document.getElementById("selectedDeck").textContent = `当前牌组: ${deckTitle}`;
        }
        await Promise.all([refreshDeckList(), refreshDeckTypeStats(), refreshCardList(), refreshUserDeck(), refreshCompletionStatus()]);
      });
    });
    row.appendChild(btn);
    container.appendChild(row);
  });
}

async function refreshCardList() {
  const keyword = document.getElementById("keyword").value.trim();
  const role = Number(document.getElementById("cardRoleFilter").value);
  const type = Number(document.getElementById("typeFilter").value);
  const cost = Number(document.getElementById("costFilter").value);
  const rarity = Number(document.getElementById("rarityFilter").value);
  const isBase = document.getElementById("upgradeFilter").value === "true";
  const summary = document.getElementById("cardSearchMeta");

  const cardsRaw = await invoke("GetCardList", keyword, type, role, cost, rarity, isBase);
  const cards = Array.isArray(cardsRaw) ? cardsRaw : [];
  const container = document.getElementById("cardList");
  if (summary) {
    summary.textContent = cards.length ? `命中 ${cards.length} 条` : "命中 0 条";
  }
  if (!cards.length) {
    showEmpty(container, "未命中卡牌");
    return;
  }

  container.innerHTML = "";
  cards.forEach((card) => {
    const row = document.createElement("div");
    row.className = "item";
    const badge = card.need ? `<span class="badge need">目标卡</span>` : `<span class="badge">非目标</span>`;
    const active = card.highlight ? `<span class="badge need">已加入</span>` : "";
    row.innerHTML = `
      <div class="item-meta">
        <img src="${toImagePath(card.image)}" alt="${card.id}" data-card-id="${card.id}" />
        <div>
          <strong>${card.id}</strong>
          <div>${badge} ${active}</div>
        </div>
      </div>
    `;
    const btn = document.createElement("button");
    btn.textContent = "加入用户卡组";
    btn.disabled = !state.selectedDeckId;
    btn.addEventListener("click", async () => {
      await runAction("加入卡牌", async () => {
        await invoke("AddCardToUserDeck", card.id);
        await Promise.all([refreshDeckTypeStats(), refreshUserDeck(), refreshCardList(), refreshCompletionStatus()]);
      });
    });
    row.appendChild(btn);
    const image = row.querySelector("img[data-card-id]");
    if (image) {
      image.addEventListener("click", () => {
        openCardDetail(card.id);
      });
    }
    container.appendChild(row);
  });
}

async function refreshUserDeck() {
  const container = document.getElementById("userDeck");
  if (!state.selectedDeckId) {
    showEmpty(container, "请先选择牌组");
    return;
  }

  const deck = await invoke("GetUserDeck");
  const cards = Array.isArray(deck?.cards) ? deck.cards : [];
  if (!cards.length) {
    showEmpty(container, "用户卡组为空");
    return;
  }

  const grouped = cards.reduce((acc, card) => {
    const key = card.id;
    if (!acc[key]) {
      acc[key] = {
        id: card.id,
        image: card.image,
        need: card.need,
        count: 0,
      };
    }
    acc[key].count += 1;
    return acc;
  }, {});
  const groupedCards = Object.values(grouped).sort((a, b) => a.id.localeCompare(b.id));

  container.innerHTML = "";
  groupedCards.forEach((card) => {
    const row = document.createElement("div");
    row.className = "item";
    row.innerHTML = `
      <div class="item-meta">
        <img src="${toImagePath(card.image)}" alt="${card.id}" data-card-id="${card.id}" />
        <div>
          <strong>${card.id}</strong>
          <div>
            ${card.need ? '<span class="badge need">目标卡</span>' : '<span class="badge">额外卡</span>'}
            <span class="badge">x${card.count}</span>
          </div>
        </div>
      </div>
    `;
    const btn = document.createElement("button");
    btn.textContent = "移除";
    btn.addEventListener("click", async () => {
      await runAction("移除卡牌", async () => {
        await invoke("RemoveCardFromUserDeck", card.id);
        await Promise.all([refreshDeckTypeStats(), refreshUserDeck(), refreshCardList(), refreshCompletionStatus()]);
      });
    });
    row.appendChild(btn);
    const image = row.querySelector("img[data-card-id]");
    if (image) {
      image.addEventListener("click", () => {
        openCardDetail(card.id);
      });
    }
    container.appendChild(row);
  });
}

function bindEvents() {
  document.getElementById("refreshDecksBtn").addEventListener("click", () => {
    runAction("刷新牌组", refreshDeckList);
  });
  document.getElementById("deckRoleFilter").addEventListener("change", () => {
    runAction("筛选牌组", refreshDeckList);
  });
  ["deckTypeKeyword", "deckTypeFilter", "deckTypeUpgradeFilter"].forEach((id) => {
    const el = document.getElementById(id);
    if (id === "deckTypeKeyword") {
      el.addEventListener("keydown", (e) => {
        if (e.key === "Enter") {
          runAction("筛选牌组统计", refreshDeckTypeStats);
        }
      });
    } else {
      el.addEventListener("change", () => {
        runAction("筛选牌组统计", refreshDeckTypeStats);
      });
    }
  });
  document.getElementById("clearDeckTypeFiltersBtn").addEventListener("click", () => {
    runAction("清空牌组统计筛选", async () => {
      resetDeckTypeFilters();
      await refreshDeckTypeStats();
    });
  });

  document.getElementById("searchCardsBtn").addEventListener("click", () => {
    runAction("搜索卡牌", refreshCardList);
  });
  document.getElementById("clearCardFiltersBtn").addEventListener("click", () => {
    runAction("清空卡牌筛选", async () => {
      resetCardFilters();
      await refreshCardList();
    });
  });
  ["cardRoleFilter", "typeFilter", "costFilter", "rarityFilter", "upgradeFilter"].forEach((id) => {
    document.getElementById(id).addEventListener("change", () => {
      runAction("筛选卡牌", refreshCardList);
    });
  });
  document.getElementById("keyword").addEventListener("keydown", (e) => {
    if (e.key === "Enter") {
      runAction("搜索卡牌", refreshCardList);
    }
  });

  document.getElementById("clearUserDeckBtn").addEventListener("click", async () => {
    await runAction("清空用户卡组", async () => {
      await invoke("ClearUserDeck");
      await Promise.all([refreshDeckTypeStats(), refreshUserDeck(), refreshCardList(), refreshCompletionStatus()]);
    });
  });

  document.getElementById("closeModalBtn").addEventListener("click", closeCardDetail);
  document.querySelector('[data-close-modal="true"]').addEventListener("click", closeCardDetail);
  document.addEventListener("keydown", (e) => {
    if (e.key === "Escape") {
      closeCardDetail();
    }
  });
}

async function init() {
  await loadDisplayLabels();
  setOptions(document.getElementById("deckRoleFilter"), ROLE_LABELS);
  setOptions(document.getElementById("deckTypeFilter"), TYPE_LABELS);
  setOptions(document.getElementById("deckTypeUpgradeFilter"), DECK_UPGRADE_FILTER_LABELS, "all");
  setOptions(document.getElementById("cardRoleFilter"), ROLE_LABELS);
  setOptions(document.getElementById("typeFilter"), TYPE_LABELS);
  setOptions(document.getElementById("costFilter"), COST_LABELS);
  setOptions(document.getElementById("rarityFilter"), RARITY_LABELS);
  setOptions(document.getElementById("upgradeFilter"), UPGRADE_LABELS, "true");
  bindEvents();

  try {
    document.getElementById("selectedDeck").textContent = "当前牌组: 未选择";
    updateCompletionLabel(false);
    await Promise.all([refreshDeckList(), refreshDeckTypeStats(), refreshCardList(), refreshUserDeck(), refreshCompletionStatus()]);
  } catch (err) {
    console.error(err);
    alert(`初始化失败: ${err.message}`);
  }
}

window.addEventListener("DOMContentLoaded", init);
