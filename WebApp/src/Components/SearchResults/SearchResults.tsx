import React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import cn from 'classnames';

import { Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle } from '@material-ui/core';

import { ResultEntity } from '../../Store/reducer';
import { searchQuerySelector, mailSelector } from '../../Store/selectors';
import { sendMail } from '../../Store/actions';

import { Button } from '../Button';

import s from './SearchResults.module.scss';

export const SearchResults = ({ results }: any) => {
  const dispatch = useDispatch();

  const searchQuery = useSelector(searchQuerySelector);
  const mailMessage = useSelector(mailSelector);
  const mailAddress = Object.keys(mailMessage)[0];
  const mailText = Object.values(mailMessage)[0];
  const [openDialog, setOpenDialog] = React.useState<boolean>(false);
  const determineReputation = (reputation: string) => {
    switch (reputation) {
      case 'Хорошая':
        return 'good';
      case 'Средняя':
        return 'medium';
      case 'Недобросовестная':
        return 'bad';
      case 'Неизвестно':
        return 'unknown';

      default:
        return 'unknown';
    }
  };

  const handleSendData = (chosenItem: ResultEntity, template = 'potato', product: string) => {
    const dataToSend = {
      template: template,
      fill_results: [chosenItem],
      product: product,
    };

    dispatch(sendMail(dataToSend));
    setOpenDialog(true);
  };

  const handleCloseDialog = () => {
    setOpenDialog(false);
  };

  return (
    <div className={s.SearchResults}>
      {results.map((item: ResultEntity, index: number) => {
        return (
          <div className={s.SearchResult} key={index}>
            <div className={s.SearchResultField}>
              <span>Название компании:</span>
              <h2>{item.company_name}</h2>
            </div>
            <div className={s.SearchResultField}>
              <span>Контактное лицо:</span>
              {item?.contact_persons?.map((contact: string) => (
                <h3>{contact}</h3>
              ))}
            </div>
            <div className={s.SearchResultField}>
              <span>E-mail:</span>
              {item?.emails?.map((email, index) => (
                <h3 key={index}>{email}</h3>
              ))}
            </div>
            <div className={s.SearchResultField}>
              <span>Телефоны: </span>
              {item?.phones?.map((phone, index) => (
                <h3 key={index}>{phone}</h3>
              ))}
            </div>
            <div className={s.SearchResultField}>
              <span>Стартовая цена: </span>
              <h3>{item.average_capitalization}</h3>
            </div>
            <div className={s.SearchResultField}>
              <span>Репутация: </span>
              <h3 className={cn(s.Reputation, s[determineReputation(item.reputation)])}>{item.reputation}</h3>
            </div>
            <Button className={s.ButtonOverride} onClick={() => handleSendData(item, 'potato', searchQuery)}>
              Отправить запрос
            </Button>
          </div>
        );
      })}
      <Dialog onClose={handleCloseDialog} open={openDialog}>
        <DialogTitle>Запрос успешно отправлен</DialogTitle>
        <DialogContent>
          <DialogContentText>Сообщение отправлено по адресу: {mailAddress}</DialogContentText>
          <DialogContentText>Текст сообщения: {mailText}</DialogContentText>
          <DialogContentText>Осталось только дождаться ответа уже на своей почте.</DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button className={s.ButtonOverride} onClick={handleCloseDialog}>
            Закрыть
          </Button>
        </DialogActions>
      </Dialog>
    </div>
  );
};
