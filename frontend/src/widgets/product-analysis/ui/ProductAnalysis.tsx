import { useState } from 'react';
import { ProductForm } from '@/entities/product/ui/ProductForm';
import { AnalysisResult } from '@/entities/product/ui/AnalysisResult';
import { ProductFormData, AnalysisResult as AnalysisResultType } from '@/entities/product/types';
import { useAnalyzeProductMutation } from '@/entities/product/api';

export const ProductAnalysis = () => {
  const [result, setResult] = useState<AnalysisResultType | null>(null);
  const [analyzeProduct, { isLoading }] = useAnalyzeProductMutation();

  const handleSubmit = async (data: ProductFormData) => {
    try {
      const response = await analyzeProduct(data).unwrap();
      setResult(response);
    } catch (error) {
      console.error('Ошибка:', error);
      const errorMessage = error instanceof Error ? error.message : 'Неизвестная ошибка';
      alert(`Произошла ошибка при анализе товара: ${errorMessage}`);
    }
  };

  return (
    <div className="min-h-screen bg-gray-100 py-6 flex flex-col justify-center sm:py-12">
      <div className="relative py-3 sm:max-w-xl sm:mx-auto">
        <div className="relative px-4 py-10 bg-white shadow-lg sm:rounded-3xl sm:p-20">
          <div className="max-w-md mx-auto">
            <div className="divide-y divide-gray-200">
              <div className="py-8 text-base leading-6 space-y-4 text-gray-700 sm:text-lg sm:leading-7">
                <h1 className="text-2xl font-bold text-center mb-8">
                  Анализ товаров
                </h1>
                <ProductForm onSubmit={handleSubmit} isLoading={isLoading} />
                <AnalysisResult result={result} />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}; 