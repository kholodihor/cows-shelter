import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { useWidth } from '@/hooks/useWidth';
import { useAppDispatch, useAppSelector } from '@/store/hook';
import { openModal } from '@/store/slices/modalSlice';
import { setActiveLink } from '@/store/slices/observationSlice';
import { useInView } from 'react-intersection-observer';

// Redux actions
import { fetchImages } from '@/store/slices/gallerySlice';

// Components
import ShareIcon from '../icons/ShareIcon';
import ShareModal from '../modals/ShareModal';
import ZoomArrow from '@/components/icons/ZoomArrow';
import Slider from '@/components/Slider';
import LightBox from './LightBox';
import LoadingSpinner from '../common/LoadingSpinner';
import ErrorMessage from '../common/ErrorMessage';

// Styles and utilities
import '@/styles/gallery.css';
import { chunkArray } from '@/utils/chunkArray';

const Gallery = () => {
  const screenWidth = useWidth();
  const { t } = useTranslation();
  const [currentPage, setCurrentPage] = useState(1);
  const [itemsPerPage, setItemsPerPage] = useState(6);
  const [image, setImage] = useState(0);
  const [showModal, setShowModal] = useState(false);
  const [pagesLength, setPagesLength] = useState(0);

  const dispatch = useAppDispatch();
  const isModalOpen = useAppSelector((state) => state.modals.isModalOpen);
  const type = useAppSelector((state) => state.modals.type);

  // Get gallery data from Redux store
  const { images, loading, error } = useAppSelector((state) => ({
    images: state.gallery.images,
    loading: state.gallery.loading,
    error: state.gallery.error
  }));

  const { ref, inView } = useInView({
    threshold: 0.5
  });

  const totalLength = images.length;

  useEffect(() => {
    if (screenWidth > 1280) {
      setItemsPerPage(6);
    }
    if (screenWidth >= 768 && screenWidth < 1280) {
      setItemsPerPage(4);
    }
    if (screenWidth > 320 && screenWidth < 768) {
      setItemsPerPage(1);
    }
  }, [screenWidth]);

  useEffect(() => {
    const pagesNumber = totalLength / itemsPerPage;
    setPagesLength(pagesNumber < 5 ? pagesNumber : 5);
  }, [totalLength, itemsPerPage]);

  useEffect(() => {
    const pagesNumber = totalLength / itemsPerPage;
    setPagesLength(pagesNumber < 5 ? pagesNumber : 5);
  }, [totalLength, itemsPerPage]);

  // Fetch gallery data on component mount
  useEffect(() => {
    dispatch(fetchImages())
      .unwrap()
      .catch((error) => {
        console.error('Failed to fetch gallery images:', error);
      });
  }, [dispatch]);

  console.log(currentPage);
  useEffect(() => {
    if (inView) {
      dispatch(setActiveLink('#gallery'));
    } else {
      dispatch(setActiveLink(''));
    }
  }, [inView, dispatch]);

  // Show loading state
  if (loading) {
    return (
      <section id="gallery" className="relative py-12">
        <div className="mx-auto px-5 sm:w-[480px] md:w-[768px] md:px-8 lg:w-[1280px] lg:px-[55px]">
          <div className="flex h-64 items-center justify-center">
            <LoadingSpinner />
          </div>
        </div>
      </section>
    );
  }

  // Show error state
  if (error) {
    return (
      <section id="gallery" className="relative py-12">
        <div className="mx-auto px-5 sm:w-[480px] md:w-[768px] md:px-8 lg:w-[1280px] lg:px-[55px]">
          <ErrorMessage
            message={error}
            onRetry={() => dispatch(fetchImages())}
          />
        </div>
      </section>
    );
  }

  // Show empty state if no images
  if (images.length === 0) {
    return (
      <section id="gallery" className="relative py-12">
        <div className="mx-auto px-5 sm:w-[480px] md:w-[768px] md:px-8 lg:w-[1280px] lg:px-[55px]">
          <p className="text-center text-lg">{t('gallery:noImages')}</p>
        </div>
      </section>
    );
  }

  // Split images into chunks for the slider based on itemsPerPage
  const imageChunks = chunkArray(images, itemsPerPage);

  return (
    <section id="gallery" ref={ref} className="relative">
      <div className="mx-auto px-5 sm:w-[480px] md:w-[768px] md:px-8 md:py-12 lg:w-[1280px] lg:px-[55px]">
        {isModalOpen && type === 'lightbox' && (
          <LightBox images={images} image={image} />
        )}

        {screenWidth >= 768 && (
          <Slider
            title={t('gallery:gallery')}
            setCurrentPage={setCurrentPage}
            pagesLength={pagesLength}
          >
            <div className="gridContainer ml-4 w-full ">
              {imageChunks &&
              imageChunks.length > 0 &&
              currentPage > 0 &&
              currentPage <= imageChunks.length ? (
                imageChunks[currentPage - 1].map((item, index) => (
                  <div
                    key={item.id}
                    className={`gridItem gridItem-- relative  h-[281px] min-w-[282px]  max-w-[456px] overflow-hidden${
                      index + 1
                    }`}
                  >
                    <img
                      src={item.image_url || '/placeholder-gallery.jpg'}
                      alt="Gallery image"
                      className="h-full w-full object-cover transition duration-500 ease-in hover:scale-110"
                      loading="lazy"
                    />

                    <div
                      onClick={() => {
                        setImage(index),
                          dispatch(openModal({ data: {}, type: 'lightbox' }));
                      }}
                      className="absolute bottom-4 left-4 z-50 cursor-pointer"
                    >
                      <ZoomArrow />
                    </div>
                  </div>
                ))
              ) : (
                <p className="text-black">Сервер не відповідає</p>
              )}
            </div>
          </Slider>
        )}
        {screenWidth < 768 && (
          <Slider
            title={t('gallery:gallery')}
            setCurrentPage={setCurrentPage}
            pagesLength={pagesLength}
          >
            <div
              className={`relative mx-auto h-[281px] w-full  min-w-[282px] overflow-hidden sm:w-[70%]`}
            >
              <img
                src={images[0]?.image_url || '/placeholder-gallery.jpg'}
                alt="Gallery image"
                className="h-full w-full object-cover transition duration-500 ease-in hover:scale-110"
                loading="lazy"
              />
              <div
                onClick={() => setShowModal(true)}
                className="absolute bottom-2 left-2 flex cursor-pointer items-center justify-center rounded-full p-2 hover:bg-[rgba(150,150,150,0.5)]"
                title="Share in Social Media"
              >
                <ShareIcon />
              </div>
              {showModal && (
                <ShareModal
                  activeImage={images[0]?.image_url || ''}
                  setShowModal={setShowModal}
                />
              )}
            </div>
          </Slider>
        )}
      </div>
    </section>
  );
};

export default Gallery;
