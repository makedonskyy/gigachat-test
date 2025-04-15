import { useForm } from 'react-hook-form';
import { ProductFormData } from '../types';
import { BeakerIcon } from '@heroicons/react/24/outline';

interface ProductFormProps {
  onSubmit: (data: ProductFormData) => void;
  isLoading: boolean;
}

export const ProductForm = ({ onSubmit, isLoading }: ProductFormProps) => {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<ProductFormData>();

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
      <div>
        <label htmlFor="name" className="block text-sm font-medium text-gray-700">
          Название товара
        </label>
        <input
          type="text"
          id="name"
          {...register('name', { required: 'Введите название товара' })}
          className="input-field"
        />
        {errors.name && (
          <p className="error-message">{errors.name.message}</p>
        )}
      </div>

      <div>
        <label htmlFor="category" className="block text-sm font-medium text-gray-700">
          Категория
        </label>
        <input
          type="text"
          id="category"
          {...register('category', { required: 'Введите категорию' })}
          className="input-field"
        />
        {errors.category && (
          <p className="error-message">{errors.category.message}</p>
        )}
      </div>

      <div>
        <label htmlFor="keywords" className="block text-sm font-medium text-gray-700">
          Ключевые слова (через запятую)
        </label>
        <input
          type="text"
          id="keywords"
          {...register('keywords', { required: 'Введите ключевые слова' })}
          className="input-field"
        />
        {errors.keywords && (
          <p className="error-message">{errors.keywords.message}</p>
        )}
      </div>

      <button
        type="submit"
        disabled={isLoading}
        className="btn-primary"
      >
        {isLoading ? (
          <>
            <BeakerIcon className="h-5 w-5 mr-2 animate-spin" />
            Анализируем...
          </>
        ) : (
          'Анализировать'
        )}
      </button>
    </form>
  );
}; 