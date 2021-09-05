import { EActionTypes } from './actionTypes';

export interface ResultEntity {
  emails: string[];
  phones: string[];
  contact_persons: string[];
  company_name: string;
  average_capitalization: string;
  reputation: string;
}

export interface ResultResponse {
  data: any;
  [key: number]: ResultEntity[];
}

export interface TaskStatusEntity {
  taskId: number;
  taskStatus: string;
}

export interface IMainStore {
  pending: boolean;
  pendingText: string;
  error: null | Error;
  results: ResultEntity[];
  taskIds: number[];
  tasksStatus: TaskStatusEntity[];
}

export const initialState: IMainStore = {
  pending: false,
  pendingText: 'Запускаем поиск...',
  error: null,
  results: [],
  taskIds: [],
  tasksStatus: [],
};

export const reducer = (state: IMainStore = initialState, action: any) => {
  const { type, payload } = action;

  switch (type) {
    case EActionTypes.LOAD_RESULTS_REQUEST: {
      return {
        ...state,
        pending: true,
      };
    }

    case EActionTypes.LOAD_RESULTS_SUCCESS: {
      return {
        ...state,
        results: payload,
        pending: false,
      };
    }

    case EActionTypes.LOAD_RESULTS_FAILURE: {
      return {
        ...state,
        pending: false,
        error: payload,
      };
    }

    case EActionTypes.SET_TASK_IDS: {
      return {
        ...state,
        taskIds: payload,
      };
    }

    case EActionTypes.SET_TASK_STATUS: {
      return {
        ...state,
        tasksStatus: payload,
      };
    }

    case EActionTypes.SEND_MAIL: {
      return {
        ...state,
        mail: payload,
      };
    }

    case EActionTypes.SET_SEARCH_QUERY: {
      return {
        ...state,
        searchQuery: payload,
      };
    }

    case EActionTypes.SET_PROGRESS_TEXT: {
      return {
        ...state,
        pendingText: payload,
      };
    }

    default:
      return state;
  }
};
