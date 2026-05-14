const fs = require("fs");
const path = require("path");
const { Client } = require("@notionhq/client");

const notion = new Client({
  auth: process.env.NOTION_TOKEN,
});

const NOTION_DATABASE_ID = process.env.NOTION_DATABASE_ID;
const GITHUB_TOKEN = process.env.GITHUB_TOKEN;
const GITHUB_REPOSITORY = process.env.GITHUB_REPOSITORY;

const AUTO_START = "[AUTO_GENERATED_FROM_README_START]";
const AUTO_END = "[AUTO_GENERATED_FROM_README_END]";
const MANUAL_START = "[MANUAL_PORTFOLIO_SECTION_START]";

if (!process.env.NOTION_TOKEN) {
  throw new Error("Missing NOTION_TOKEN");
}

if (!NOTION_DATABASE_ID) {
  throw new Error("Missing NOTION_DATABASE_ID");
}

if (!GITHUB_REPOSITORY) {
  throw new Error("Missing GITHUB_REPOSITORY");
}

function plainTextFromRichText(richText = []) {
  return richText.map((item) => item.plain_text || "").join("");
}

function richText(content) {
  const text = String(content || "");

  return [
    {
      type: "text",
      text: {
        content: text.slice(0, 2000),
      },
    },
  ];
}

function paragraph(content) {
  return {
    object: "block",
    type: "paragraph",
    paragraph: {
      rich_text: richText(content),
    },
  };
}

function heading1(content) {
  return {
    object: "block",
    type: "heading_1",
    heading_1: {
      rich_text: richText(content),
    },
  };
}

function heading2(content) {
  return {
    object: "block",
    type: "heading_2",
    heading_2: {
      rich_text: richText(content),
    },
  };
}

function heading3(content) {
  return {
    object: "block",
    type: "heading_3",
    heading_3: {
      rich_text: richText(content),
    },
  };
}

function bullet(content) {
  return {
    object: "block",
    type: "bulleted_list_item",
    bulleted_list_item: {
      rich_text: richText(content),
    },
  };
}

function codeBlock(content, language = "plain text") {
  return {
    object: "block",
    type: "code",
    code: {
      language,
      rich_text: richText(content),
    },
  };
}

function divider() {
  return {
    object: "block",
    type: "divider",
    divider: {},
  };
}

function normalizeVisibility(value) {
  return value === "private" ? "Private" : "Public";
}

function normalizeStatus(repo) {
  if (repo.archived) return "보관";
  if (repo.disabled) return "중단";
  return "운영 중";
}

function inferRole(language) {
  const lang = String(language || "").toLowerCase();

  if (["javascript", "typescript", "html", "css"].includes(lang)) {
    return "Frontend";
  }

  if (
    ["go", "python", "java", "kotlin", "rust", "c", "c++", "c#"].includes(lang)
  ) {
    return "Backend";
  }

  return "Project";
}

function getReadme() {
  const readmePath = path.join(process.cwd(), "README.md");

  if (!fs.existsSync(readmePath)) {
    return "";
  }

  return fs.readFileSync(readmePath, "utf8");
}

function extractReadmeSummary(readme, repo) {
  const lines = readme.split(/\r?\n/).map((line) => line.trim());

  const firstTitle = lines.find((line) => line.startsWith("# "));
  const title = firstTitle ? firstTitle.replace(/^#\s+/, "").trim() : repo.name;

  const quote = lines
    .filter((line) => line.startsWith(">"))
    .map((line) => line.replace(/^>\s?/, "").trim())
    .find(Boolean);

  const firstParagraph = lines.find((line) => {
    if (!line) return false;
    if (line.startsWith("#")) return false;
    if (line.startsWith("-")) return false;
    if (line.startsWith("```")) return false;
    if (line.startsWith(">")) return false;
    return true;
  });

  const summary =
    quote ||
    repo.description ||
    firstParagraph ||
    `${repo.name} 저장소의 README 기반 포트폴리오 초안입니다.`;

  return {
    title,
    summary,
  };
}

function markdownToBlocks(markdown) {
  const lines = markdown.split(/\r?\n/);
  const blocks = [];
  let inCode = false;
  let codeLines = [];
  let codeLanguage = "plain text";

  for (const rawLine of lines) {
    const line = rawLine.trimEnd();

    if (line.startsWith("```")) {
      if (!inCode) {
        inCode = true;
        codeLanguage = line.replace(/^```/, "").trim() || "plain text";
        codeLines = [];
      } else {
        inCode = false;
        blocks.push(
          codeBlock(codeLines.join("\n").slice(0, 2000), codeLanguage),
        );
        codeLines = [];
        codeLanguage = "plain text";
      }
      continue;
    }

    if (inCode) {
      codeLines.push(rawLine);
      continue;
    }

    if (!line.trim()) {
      continue;
    }

    if (line.startsWith("# ")) {
      blocks.push(heading1(line.replace(/^#\s+/, "")));
      continue;
    }

    if (line.startsWith("## ")) {
      blocks.push(heading2(line.replace(/^##\s+/, "")));
      continue;
    }

    if (line.startsWith("### ")) {
      blocks.push(heading3(line.replace(/^###\s+/, "")));
      continue;
    }

    if (line.startsWith("- ")) {
      blocks.push(bullet(line.replace(/^-\s+/, "")));
      continue;
    }

    if (line.startsWith("* ")) {
      blocks.push(bullet(line.replace(/^\*\s+/, "")));
      continue;
    }

    if (line.startsWith(">")) {
      blocks.push(paragraph(line.replace(/^>\s?/, "")));
      continue;
    }

    blocks.push(paragraph(line));
  }

  return blocks.slice(0, 80);
}

async function fetchGithubRepo() {
  const response = await fetch(
    `https://api.github.com/repos/${GITHUB_REPOSITORY}`,
    {
      headers: {
        Authorization: `Bearer ${GITHUB_TOKEN}`,
        Accept: "application/vnd.github+json",
        "User-Agent": "readme-notion-sync",
      },
    },
  );

  if (!response.ok) {
    const body = await response.text();
    throw new Error(`GitHub repo fetch failed: ${response.status} ${body}`);
  }

  return response.json();
}

async function getDatabaseTitlePropertyName() {
  const database = await notion.databases.retrieve({
    database_id: NOTION_DATABASE_ID,
  });

  const properties = database.properties;

  for (const [name, property] of Object.entries(properties)) {
    if (property.type === "title") {
      return name;
    }
  }

  throw new Error("Notion database title property not found");
}

async function findPageByTitle(titlePropertyName, title) {
  const response = await notion.databases.query({
    database_id: NOTION_DATABASE_ID,
    filter: {
      property: titlePropertyName,
      title: {
        equals: title,
      },
    },
  });

  return response.results[0] || null;
}

function makeProperties(titlePropertyName, repo) {
  const language = repo.language || "Unknown";

  const properties = {
    [titlePropertyName]: {
      title: richText(repo.name),
    },
  };

  properties["공개여부"] = {
    select: {
      name: normalizeVisibility(repo.visibility),
    },
  };

  properties["기간"] = {
    date: {
      start: repo.created_at ? repo.created_at.slice(0, 10) : null,
    },
  };

  properties["링크"] = {
    url: repo.html_url || null,
  };

  properties["상태"] = {
    select: {
      name: normalizeStatus(repo),
    },
  };

  properties["언어"] = {
    multi_select: {
      name: language,
    },
  };

  properties["역할"] = {
    select: {
      name: inferRole(language),
    },
  };

  return properties;
}

function buildAutoBlocks(repo, readme) {
  const summary = extractReadmeSummary(readme, repo);

  return [
    paragraph(AUTO_START),
    heading1(`${repo.name}`),
    paragraph(summary.summary),
    divider(),

    heading2("개요"),
    paragraph(`저장소: ${repo.full_name}`),
    paragraph(`설명: ${repo.description || "-"}`),
    paragraph(`주요 언어: ${repo.language || "-"}`),
    paragraph(`공개 여부: ${normalizeVisibility(repo.visibility)}`),
    paragraph(
      `생성일: ${repo.created_at ? repo.created_at.slice(0, 10) : "-"}`,
    ),
    paragraph(
      `최근 수정일: ${repo.updated_at ? repo.updated_at.slice(0, 10) : "-"}`,
    ),

    heading2("README 기반 초안"),
    paragraph("아래 영역은 GitHub README.md를 기준으로 자동 생성됩니다."),

    ...markdownToBlocks(readme),

    divider(),
    paragraph(AUTO_END),
  ];
}

function buildManualBlocks() {
  return [
    heading2("수동 작성 영역"),
    paragraph(MANUAL_START),
    paragraph(
      "이 아래에는 Notion에서 직접 포트폴리오 내용을 보강하세요. README가 수정되어도 이 영역은 자동화 대상에서 제외하는 것을 목표로 합니다.",
    ),

    heading2("왜 만들었나 — 문제 정의"),
    paragraph("작성 필요"),

    heading2("핵심 구현 — 기술적 의사결정"),
    paragraph("작성 필요"),

    heading2("주요 성과"),
    paragraph("작성 필요"),

    heading2("회고"),
    paragraph("작성 필요"),
  ];
}

async function listChildren(blockId) {
  const results = [];
  let cursor = undefined;

  do {
    const response = await notion.blocks.children.list({
      block_id: blockId,
      start_cursor: cursor,
      page_size: 100,
    });

    results.push(...response.results);
    cursor = response.has_more ? response.next_cursor : undefined;
  } while (cursor);

  return results;
}

function getBlockPlainText(block) {
  const type = block.type;
  const value = block[type];

  if (!value || !value.rich_text) {
    return "";
  }

  return plainTextFromRichText(value.rich_text);
}

async function archiveBlocks(blocks) {
  for (const block of blocks) {
    await notion.blocks.update({
      block_id: block.id,
      archived: true,
    });
  }
}

function chunkArray(array, size) {
  const chunks = [];

  for (let i = 0; i < array.length; i += size) {
    chunks.push(array.slice(i, i + size));
  }

  return chunks;
}

async function appendBlocks(blockId, blocks, after) {
  const chunks = chunkArray(blocks, 80);
  let currentAfter = after;

  for (const chunk of chunks) {
    const payload = {
      block_id: blockId,
      children: chunk,
    };

    if (currentAfter) {
      payload.after = currentAfter;
    }

    const response = await notion.blocks.children.append(payload);
    const appended = response.results || [];

    if (appended.length > 0) {
      currentAfter = appended[appended.length - 1].id;
    }
  }
}

async function replaceAutoGeneratedArea(pageId, autoBlocks) {
  const children = await listChildren(pageId);

  const startIndex = children.findIndex(
    (block) => getBlockPlainText(block) === AUTO_START,
  );
  const endIndex = children.findIndex(
    (block) => getBlockPlainText(block) === AUTO_END,
  );

  if (startIndex === -1 || endIndex === -1 || endIndex <= startIndex) {
    await appendBlocks(pageId, [...autoBlocks, ...buildManualBlocks()]);
    return;
  }

  const startBlock = children[startIndex];
  const blocksToDelete = children.slice(startIndex + 1, endIndex + 1);

  await archiveBlocks(blocksToDelete);

  const blocksWithoutStartMarker = autoBlocks.slice(1);
  await appendBlocks(pageId, blocksWithoutStartMarker, startBlock.id);
}

async function main() {
  const repo = await fetchGithubRepo();
  const readme = getReadme();

  const titlePropertyName = await getDatabaseTitlePropertyName();
  const properties = makeProperties(titlePropertyName, repo);
  const existingPage = await findPageByTitle(titlePropertyName, repo.name);
  const autoBlocks = buildAutoBlocks(repo, readme);

  if (!existingPage) {
    const page = await notion.pages.create({
      parent: {
        database_id: NOTION_DATABASE_ID,
      },
      properties,
      children: [...autoBlocks, ...buildManualBlocks()],
    });

    console.log(`Created Notion page: ${page.id}`);
    return;
  }

  await notion.pages.update({
    page_id: existingPage.id,
    properties,
  });

  await replaceAutoGeneratedArea(existingPage.id, autoBlocks);

  console.log(`Updated Notion page: ${existingPage.id}`);
}

main().catch((error) => {
  console.error(error);
  process.exit(1);
});
