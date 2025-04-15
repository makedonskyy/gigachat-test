# Product Analysis Frontend

Фронтенд-часть приложения для анализа товаров с использованием GigaChat API.

## Технологии

- React 18
- TypeScript
- Vite
- Redux Toolkit + RTK Query
- Tailwind CSS
- Feature-Sliced Design

## Структура проекта

```
src/
  ├── app/          # Инициализация приложения
  ├── entities/     # Бизнес-сущности
  ├── features/     # Функциональность
  ├── shared/       # Переиспользуемый код
  └── widgets/      # Композиционные компоненты
```

## Установка и запуск

1. Установите зависимости:
```bash
npm install
```

2. Запустите приложение в режиме разработки:
```bash
npm run dev
```

3. Для сборки проекта:
```bash
npm run build
```

## Особенности

- Feature-Sliced Design архитектура для лучшей масштабируемости
- RTK Query для управления состоянием и API запросами
- Tailwind CSS для стилизации компонентов
- TypeScript для типобезопасности