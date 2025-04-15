import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';
import { ProductFormData, AnalysisResult } from '@/entities/product';

export const api = createApi({
  reducerPath: 'api',
  baseQuery: fetchBaseQuery({ baseUrl: 'http://localhost:8080' }),
  endpoints: () => ({}),
});

export const { middleware: apiMiddleware, reducer: apiReducer } = api; 