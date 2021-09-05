import { EActionTypes } from './actionTypes';
import { ResultEntity, ResultResponse } from './reducer';
import api from '../Api';

export const startSearch = (searchData: any) => async (dispatch: any) => {
  dispatch({ type: EActionTypes.LOAD_RESULTS_REQUEST });
  try {
    const taskId = await api.search(searchData);
    dispatch({ type: EActionTypes.SET_TASK_IDS, payload: taskId.ids });
  } catch (err) {
    console.error(err);
    dispatch({ type: EActionTypes.LOAD_RESULTS_FAILURE, payload: err });
  }
};

export const fetchTaskStatus = () => async (dispatch: any, getState: any) => {
  try {
    const taskIds = getState().taskIds;
    const res = await api.getTaskStatus(taskIds);
    console.log(res);
    dispatch({ type: EActionTypes.SET_PROGRESS_TEXT, payload: res[0].progress });

    if (typeof res !== 'string') {
      if (res.every((item: any) => item.status === 'done')) {
        dispatch({ type: EActionTypes.SET_TASK_STATUS, payload: res });
      }
    }
  } catch (err) {
    console.error(err);
    dispatch({ type: EActionTypes.LOAD_RESULTS_FAILURE, payload: err });
  }
};

export const setSearchResults = () => async (dispatch: any, getState: any) => {
  dispatch({ type: EActionTypes.LOAD_RESULTS_REQUEST });
  const taskIds = getState().taskIds;

  try {
    const results: ResultResponse = await api.getTaskResults(taskIds);
    let normalizedResults: ResultEntity[] = [];

    for (let i = 0; i < taskIds.length; i++) {
      normalizedResults = normalizedResults.concat(results.data[taskIds[i]]);
    }
    dispatch({ type: EActionTypes.LOAD_RESULTS_SUCCESS, payload: normalizedResults });
  } catch (err) {
    console.error(err);
    dispatch({ type: EActionTypes.LOAD_RESULTS_FAILURE, payload: err });
  }
};

export const sendMail = (data: any) => async (dispatch: any) => {
  try {
    const res = await api.sendMail(data);
    dispatch({ type: EActionTypes.SEND_MAIL, payload: res });
  } catch (err) {
    console.error(err);
  }
};
