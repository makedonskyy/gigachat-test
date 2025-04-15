import ReactMarkdown from 'react-markdown';
import { AnalysisResult as AnalysisResultType } from '../types';

interface AnalysisResultProps {
  result: AnalysisResultType | null;
}

export const AnalysisResult = ({ result }: AnalysisResultProps) => {
  if (!result) return null;

  return (
    <div className="mt-8 bg-white shadow rounded-lg p-6">
      <h2 className="text-lg font-medium text-gray-900 mb-4">
        Результаты анализа
      </h2>
      <div className="prose prose-primary max-w-none">
        <ReactMarkdown>{result.content}</ReactMarkdown>
      </div>
    </div>
  );
}; 