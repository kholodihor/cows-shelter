import { NewsFormInput } from '@/types';
import {
  AnyAction,
  createAsyncThunk,
  createSlice,
  PayloadAction
} from '@reduxjs/toolkit';
import { AxiosError } from 'axios';
import axiosInstance from '@/utils/axios';
import { transformMinioUrlsInData } from '@/utils/minioUrlHelper';

export type Post = {
  id: string;
  title_ua: string;
  title_en: string;
  subtitle_ua: string;
  subtitle_en: string;
  content_ua: string;
  content_en: string;
  image_url: string;
  createdAt: string;
};

type ResponseWithPagination = {
  posts: Post[];
  totalLength: number;
};

type NewsState = {
  posts: Post[];
  loading: boolean;
  error: string | null;
  paginatedData: ResponseWithPagination;
};

const initialState: NewsState = {
  posts: [],
  loading: false,
  error: null,
  paginatedData: {
    posts: [],
    totalLength: 0
  }
};

export const fetchPosts = createAsyncThunk(
  'news/fetchPosts',
  async (_, { rejectWithValue }) => {
    try {
      const response = await axiosInstance.get<Post[]>('/news');
      const data = response.data;
      // Transform MinIO URLs in the response data
      const transformedData = transformMinioUrlsInData(
        Array.isArray(data) ? data : []
      );
      return transformedData;
    } catch (error) {
      const err = error as AxiosError;
      console.error('Failed to fetch news posts:', err.message);
      return rejectWithValue('Failed to load news posts');
    }
  }
);

export const fetchPostById = createAsyncThunk(
  'news/fetchPostById',
  async (id: string, { rejectWithValue }) => {
    try {
      const response = await axiosInstance.get<Post>(`/news/${id}`);
      // Transform MinIO URLs in the response data
      const transformedData = transformMinioUrlsInData(response.data);
      return transformedData;
    } catch (error) {
      const err = error as AxiosError;
      console.error(`Failed to fetch news post with id ${id}:`, err.message);
      return rejectWithValue('Failed to load news post');
    }
  }
);

export const fetchNewsWithPagination = createAsyncThunk(
  'news/fetchNewsWithPagination',
  async (query: { page: number; limit: number }, { rejectWithValue }) => {
    try {
      const response = await axiosInstance.get<{
        data: Post[];
        total: number;
      }>(`/news?page=${query.page}&limit=${query.limit}`);

      // Transform the response to match the expected format
      const transformedData: ResponseWithPagination = {
        posts: transformMinioUrlsInData(response.data?.data || []),
        totalLength: response.data?.total || 0
      };

      return transformedData;
    } catch (error) {
      const err = error as AxiosError;
      console.error('Failed to fetch paginated news:', err.message, error);
      return rejectWithValue(err.response?.data || 'Failed to load news');
    }
  }
);

export const removePost = createAsyncThunk(
  'news/removePost',
  async (id: string, { rejectWithValue }) => {
    try {
      await axiosInstance.delete(`/news/${id}`);
      return id; // Return the deleted ID for state updates
    } catch (error) {
      const err = error as AxiosError;
      console.error('Failed to delete news:', err.message, error);
      return rejectWithValue(err.response?.data || 'Failed to delete news');
    }
  }
);

// Helper function to convert file to base64
const fileToBase64 = (file: File): Promise<string> => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.readAsDataURL(file);
    reader.onload = () => resolve(reader.result as string);
    reader.onerror = (error) => reject(error);
  });
};

export const addNewPost = createAsyncThunk(
  'news/addNewPost',
  async (values: NewsFormInput, { rejectWithValue }) => {
    try {
      // Convert image to base64 if it exists
      let imageData = '';

      if (values.image && values.image[0] && values.image[0] instanceof File) {
        try {
          const file = values.image[0];
          // Check if file size is reasonable (e.g., less than 5MB)
          const MAX_FILE_SIZE = 5 * 1024 * 1024; // 5MB
          if (file.size > MAX_FILE_SIZE) {
            return rejectWithValue('Image size should be less than 5MB');
          }

          imageData = await fileToBase64(file);
        } catch (uploadError) {
          console.error('Image processing failed:', uploadError);
          return rejectWithValue('Failed to process image');
        }
      }

      // Prepare the post data with base64 image
      const postData = {
        title_ua: values.titleUa || '',
        title_en: values.titleEn || values.titleUa || '', // Fallback to UA title if EN is empty
        subtitle_ua: values.subTitleUa || '',
        subtitle_en: values.subTitleEn || values.subTitleUa || '',
        content_ua: values.contentUa || '',
        content_en: values.contentEn || values.contentUa || '',
        image_data: imageData || undefined // Only include if there's an image
      };

      console.log('Sending post data with image:', !!imageData);

      // Create the post with JSON data
      const response = await axiosInstance.post<Post>('/news', postData, {
        headers: {
          'Content-Type': 'application/json'
        }
      });

      if (!response.data) {
        throw new Error('No data received from server');
      }

      return response.data;
    } catch (error) {
      const err = error as AxiosError;
      console.error('Failed to add new post:', err.message);
      return rejectWithValue('Failed to add new post');
    }
  }
);

export const editPost = createAsyncThunk(
  'news/editPost',
  async (
    postData: { id?: string; values: NewsFormInput },
    { rejectWithValue }
  ) => {
    try {
      if (!postData.id) {
        throw new Error('Post ID is required for updating');
      }

      // Convert image to base64 if a new image is provided
      let imageData = '';

      if (
        postData.values.image?.[0] &&
        postData.values.image[0] instanceof File
      ) {
        // Check if this is a real file with actual content (not a dummy file)
        // Dummy files created in the edit form have size 0
        if (postData.values.image[0].size > 0) {
          try {
            const file = postData.values.image[0];
            // Check if file size is reasonable (e.g., less than 5MB)
            const MAX_FILE_SIZE = 5 * 1024 * 1024; // 5MB
            if (file.size > MAX_FILE_SIZE) {
              return rejectWithValue('Image size should be less than 5MB');
            }

            imageData = await fileToBase64(file);
          } catch (uploadError) {
            console.error('Image processing failed:', uploadError);
            return rejectWithValue('Failed to process image');
          }
        }
      }

      // Prepare the update data with the same structure as addNewPost
      const updateData: Record<string, any> = {
        title_ua: postData.values.titleUa || '',
        title_en: postData.values.titleEn || postData.values.titleUa || '',
        subtitle_ua: postData.values.subTitleUa || '',
        subtitle_en:
          postData.values.subTitleEn || postData.values.subTitleUa || '',
        content_ua: postData.values.contentUa || '',
        content_en: postData.values.contentEn || postData.values.contentUa || ''
      };

      // Only include image_data if a new image was provided
      // This ensures we keep the existing Cloudinary URL if no new image is uploaded
      if (imageData) {
        updateData.image_data = imageData;
      }

      console.log('Updating post with data:', updateData);

      // Update the post with JSON data
      const response = await axiosInstance.put<Post>(
        `/news/${postData.id}`,
        updateData,
        {
          headers: {
            'Content-Type': 'application/json'
          }
        }
      );

      if (!response.data) {
        throw new Error('No data received from server');
      }

      return response.data;
    } catch (err) {
      console.error('Error updating post:', err);
      return rejectWithValue(
        err instanceof Error ? err.message : 'Failed to update post'
      );
    }
  }
);

const newsSlice = createSlice({
  name: 'news',
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchPosts.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchPosts.fulfilled, (state, action) => {
        state.loading = false;
        if (Array.isArray(action.payload)) {
          state.posts = action.payload;
        } else {
          state.error = 'Invalid response format';
          state.posts = [];
        }
      })
      .addCase(fetchPostById.pending, (state) => {
        state.loading = true;
        state.error = null;
        state.posts = [];
      })
      .addCase(fetchPostById.fulfilled, (state, action) => {
        state.loading = false;
        if (
          action.payload &&
          typeof action.payload === 'object' &&
          'id' in action.payload
        ) {
          state.posts = [action.payload as Post];
        } else {
          state.error = 'Invalid post data';
        }
      })
      .addCase(fetchNewsWithPagination.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchNewsWithPagination.fulfilled, (state, action) => {
        state.loading = false;
        if (
          action.payload &&
          'posts' in action.payload &&
          'totalLength' in action.payload
        ) {
          state.paginatedData = action.payload;
        } else {
          state.error = 'Invalid pagination data';
          state.paginatedData = { posts: [], totalLength: 0 };
        }
      })
      .addCase(removePost.fulfilled, (state, action) => {
        state.posts = state.posts.filter(
          (item) => item.id !== (action.meta.arg as string)
        );
      })
      .addMatcher(isError, (state, action: PayloadAction<string>) => {
        state.loading = false;
        state.error = action.payload || 'An error occurred';
      });
  }
});

export default newsSlice.reducer;

function isError(action: AnyAction) {
  return action.type.endsWith('rejected');
}
