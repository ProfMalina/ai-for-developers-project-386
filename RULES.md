# Правила разработки проекта

## 🚫 Запрет на глобальную установку зависимостей

### Правило
**ЗАПРЕЩЕНО** устанавливать зависимости глобально. Все зависимости должны устанавливаться локально в рамках проекта.

### Обоснование
- Глобальные зависимости создают конфликты версий между проектами
- Нарушают воспроизводимость сборки на разных машинах
- Усложняют CI/CD процессы
- Делают проект зависимым от конфигурации разработчика

### Примеры нарушений

❌ **НЕЛЬЗЯ** — глобальная установка через `npm install -g`:
```bash
npm install -g typescript
npm install -g @typespec/compiler
npm install -g nodemon
```

❌ **НЕЛЬЗЯ** — использование глобальных команд вместо npx:
```bash
tsc main.ts
tsp compile .
nodemon server.js
```

### Правильное использование

✅ **МОЖНО** — использование `npx` для запуска CLI инструментов:
```bash
npx tsc --version
npx tsp compile .
npx nodemon server.js
```

✅ **МОЖНО** — использование скриптов из `package.json`:
```json
{
  "scripts": {
    "compile": "tsp compile .",
    "build": "tsc",
    "dev": "nodemon server.js"
  }
}
```

Затем запускать через:
```bash
npm run compile
npm run build
npm run dev
```

✅ **МОЖНО** — локальная установка зависимостей:
```bash
npm install --save-dev typescript
npm install --save-dev @typespec/compiler
```

### Применяется к инструментам
Это правило распространяется на все CLI инструменты, включая но не ограничиваясь:
- `tsc` → использовать `npx tsc`
- `tsp` → использовать `npx tsp`
- `nodemon` → использовать `npx nodemon`
- `eslint` → использовать `npx eslint`
- `prettier` → использовать `npx prettier`
- `jest` → использовать `npx jest`
- `webpack` → использовать `npx webpack`
- `typescript`, `@typespec/*` и другие пакеты

### Исключения
Нет исключений из этого правила. Все зависимости должны быть локальными.
