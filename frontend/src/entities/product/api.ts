import { api } from '@/shared/api/api';
import { ProductFormData, AnalysisResult } from './types';

export const productApi = api.injectEndpoints({
  endpoints: (builder) => ({
    analyzeProduct: builder.mutation<AnalysisResult, ProductFormData>({
      query: (data) => ({
        url: '/analyze',
        method: 'POST',
        body: data,
      }),
    }),
  }),
});

export const { useAnalyzeProductMutation } = productApi; 