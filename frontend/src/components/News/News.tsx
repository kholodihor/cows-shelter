import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useInView } from 'react-intersection-observer';
import { SwiperSlide } from 'swiper/react';
import { useWidth } from '@/hooks/useWidth';
import { useAppDispatch, useAppSelector } from '@/store/hook';
import { setActiveLink } from '@/store/slices/observationSlice';
import { fetchNewsWithPagination, fetchPosts } from '@/store/slices/newsSlice';
import NewsModal from '@/components/modals/NewsModal';
import Slider from '../Slider';
import NewsBlock from './NewsBlock';
import LoadingSpinner from '../common/LoadingSpinner';
import ErrorMessage from '../common/ErrorMessage';

import 'swiper/css/pagination';
import 'swiper/css';

const News = () => {
  const screenWidth = useWidth();
  const { t } = useTranslation();
  const dispatch = useAppDispatch();
  const [, setShowModal] = useState(false);
  const [currentPage, setCurrentPage] = useState(1);
  const [itemsPerPage, setItemsPerPage] = useState(5);
  const [pagesLength, setPagesLength] = useState(0);

  // Get data from Redux store
  const { posts, loading, error, paginatedData } = useAppSelector((state) => state.posts);
  const type = useAppSelector((state) => state.modals.type);
  const isModalOpen = useAppSelector((state) => state.modals.isModalOpen);

  const totalLength = paginatedData.totalLength || posts.length;
  const displayedPosts = paginatedData.posts.length > 0 ? paginatedData.posts : posts;

  const { inView } = useInView({
    threshold: 0.5,
    triggerOnce: true
  });

  // Fetch posts on component mount and when page changes
  useEffect(() => {
    const fetchData = async () => {
      try {
        if (itemsPerPage > 0) {
          const resultAction = await dispatch(fetchNewsWithPagination({ 
            page: currentPage, 
            limit: itemsPerPage 
          }));
          
          if (fetchNewsWithPagination.rejected.match(resultAction)) {
            console.error('Failed to fetch paginated news:', resultAction.error);
          }
        } else {
          await dispatch(fetchPosts());
        }
      } catch (err) {
        console.error('Failed to fetch news:', err);
      }
    };

    fetchData();
  }, [dispatch, currentPage, itemsPerPage]);

  // Handle responsive items per page
  useEffect(() => {
    if (screenWidth >= 1280) {
      setItemsPerPage(5);
    } else if (screenWidth >= 768) {
      setItemsPerPage(3);
    } else {
      setItemsPerPage(1);
    }
  }, [screenWidth]);

  // Calculate pagination
  useEffect(() => {
    const pagesNumber = Math.ceil(totalLength / itemsPerPage);
    setPagesLength(Math.min(pagesNumber, 5));
  }, [totalLength, itemsPerPage]);

  // Handle retry on error
  const handleRetry = async () => {
    try {
      if (itemsPerPage > 0) {
        const resultAction = await dispatch(fetchNewsWithPagination({ 
          page: currentPage, 
          limit: itemsPerPage 
        }));
        
        if (fetchNewsWithPagination.rejected.match(resultAction)) {
          console.error('Retry failed:', resultAction.error);
        }
      } else {
        await dispatch(fetchPosts());
      }
    } catch (err) {
      console.error('Retry error:', err);
    }
  };


  useEffect(() => {
    if (isModalOpen) {
      setShowModal(true);
    } else {
      setShowModal(false);
    }
  }, [isModalOpen]);

  useEffect(() => {
    if (inView) {
      dispatch(setActiveLink('#news'));
    } else {
      dispatch(setActiveLink(''));
    }
  }, [inView, dispatch]);

  // Loading state
  if (loading && displayedPosts.length === 0) {
    return (
      <section className="bg-white py-20 md:py-40" id="news">
        <div className="container">
          <h2 className="text-2xl font-bold md:text-4xl text-main-text-color mb-10">
            {t('news.title')}
          </h2>
          <div className="flex items-center justify-center h-64">
            <LoadingSpinner size={30} />
          </div>
        </div>
      </section>
    );
  }

  // Error state
  if (error) {
    return (
      <section className="bg-white py-20 md:py-40" id="news">
        <div className="container">
          <h2 className="text-2xl font-bold md:text-4xl text-main-text-color mb-10">
            {t('news.title')}
          </h2>
          <ErrorMessage
            message={t('news.errorLoading') || 'Failed to load news. Please try again.'}
            onRetry={handleRetry}
          />
        </div>
      </section>
    );
  }

  // Empty state
  if (displayedPosts.length === 0) {
    return (
      <section className="bg-white py-20 md:py-40" id="news">
        <div className="container">
          <h2 className="text-2xl font-bold md:text-4xl text-main-text-color mb-10">
            {t('news.title')}
          </h2>
          <p className="text-center text-gray-500">
            {t('news.noNews') || 'No news available at the moment.'}
          </p>
        </div>
      </section>
    );
  }

  // Success state
  return (
    <section className="bg-white pb-10 pt-20 md:pt-40" id="news">
      <div className="container">
        <div className="flex items-center justify-between mb-10">
          <h2 className="text-2xl font-bold md:text-4xl text-main-text-color">
            {t('news.title')}
          </h2>
          <span className="text-sm font-medium text-main-text-color">
            {t('news.viewAll')}
          </span>
        </div>
        <div className="relative">
          <div className="relative flex flex-col items-center justify-center h-full">
            {screenWidth > 767 && (
              <button
                className="absolute left-0 z-10 flex items-center justify-center w-10 h-10 -translate-y-1/2 rounded-full top-1/2 bg-main-bg-color hover:bg-gray-100 transition-colors"
                onClick={() => setCurrentPage(prev => Math.max(1, prev - 1))}
                disabled={currentPage === 1}
              >
                <svg
                  width="24"
                  height="24"
                  viewBox="0 0 24 24"
                  fill="none"
                  xmlns="http://www.w3.org/2000/svg"
                  className={`${currentPage === 1 ? 'opacity-30' : 'cursor-pointer'}`}
                >
                  <path
                    d="M15 18L9 12L15 6"
                    stroke="#1E1E1E"
                    strokeWidth="2"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  />
                </svg>
              </button>
            )}
            <Slider
              title={t('news.title')}
              pagesLength={pagesLength}
              setCurrentPage={setCurrentPage}
              isReviews={false}
              isExcursions={false}
              isPartners={false}
            >
              {displayedPosts.map((post) => (
                <SwiperSlide key={post.id}>
                  <NewsBlock post={post} onClick={() => setShowModal(true)} />
                </SwiperSlide>
              ))}
            </Slider>
            {screenWidth > 767 && (
              <button
                className="absolute right-0 z-10 flex items-center justify-center w-10 h-10 -translate-y-1/2 rounded-full top-1/2 bg-main-bg-color hover:bg-gray-100 transition-colors"
                onClick={() => setCurrentPage(prev => Math.min(pagesLength, prev + 1))}
                disabled={currentPage >= pagesLength}
              >
                <svg
                  width="24"
                  height="24"
                  viewBox="0 0 24 24"
                  fill="none"
                  xmlns="http://www.w3.org/2000/svg"
                  className={`${currentPage >= pagesLength ? 'opacity-30' : 'cursor-pointer'}`}
                >
                  <path
                    d="M9 18L15 12L9 6"
                    stroke="#1E1E1E"
                    strokeWidth="2"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  />
                </svg>
              </button>
            )}
          </div>

          {/* Pagination dots */}
          <div className="flex justify-center mt-6 space-x-2">
            {Array.from({ length: pagesLength }, (_, i) => i + 1).map((page) => (
              <button
                key={page}
                onClick={() => setCurrentPage(page)}
                className={`w-3 h-3 rounded-full transition-colors ${currentPage === page ? 'bg-main-text-color' : 'bg-gray-300'
                  }`}
                aria-label={`Go to page ${page}`}
              />
            ))}
          </div>
        </div>
      </div>
      {isModalOpen && type === 'news' && (
        <NewsModal isOpen={isModalOpen} setShowModal={setShowModal} />
      )}
    </section>
  );
};

export default News;
