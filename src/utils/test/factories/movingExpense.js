import moment from 'moment';
import { v4 } from 'uuid';

import { expenseTypes } from 'constants/ppmExpenseTypes';
import { createDocumentWithoutUploads } from 'utils/test/factories/document';
import createUpload from 'utils/test/factories/upload';

const createBaseMovingExpense = ({ serviceMemberId, creationDate = new Date() } = {}, fieldOverrides = {}) => {
  const createdAtDate = creationDate.toISOString();

  const smId = serviceMemberId || v4();
  const document = createDocumentWithoutUploads({ serviceMemberId: smId });

  return {
    id: v4(),
    ppmShipmentId: v4(),
    documentId: document.id,
    document,
    movingExpenseType: null,
    description: null,
    paidWithGtcc: null,
    amount: null,
    missingReceipt: null,
    status: null,
    reason: null,
    sitStartDate: null,
    sitEndDate: null,
    createdAt: createdAtDate,
    updatedAt: createdAtDate,
    eTag: window.btoa(createdAtDate),
    ...fieldOverrides,
  };
};

const createCompleteMovingExpense = ({ serviceMemberId, creationDate } = {}, fieldOverrides = {}) => {
  const fullFieldOverrides = {
    movingExpenseType: expenseTypes.PACKING_MATERIALS,
    description: 'Medium and large boxes',
    paidWithGtcc: false,
    amount: '7500',
    missingReceipt: false,
    ...fieldOverrides,
  };

  const movingExpense = createBaseMovingExpense({ serviceMemberId, creationDate }, fullFieldOverrides);

  if (movingExpense.document.uploads.length === 0) {
    movingExpense.document.uploads.push(createUpload({ fileName: 'expense.pdf' }));
  }

  if (movingExpense.createdAt === movingExpense.updatedAt) {
    const updatedAt = moment(movingExpense.createdAt).add(1, 'hour').toISOString();

    movingExpense.updatedAt = updatedAt;
    movingExpense.eTag = window.btoa(updatedAt);
  }

  return movingExpense;
};

const createCompleteSITMovingExpense = ({ serviceMemberId, creationDate } = {}, fieldOverrides = {}) => {
  const fullFieldOverrides = {
    movingExpenseType: expenseTypes.STORAGE,
    description: 'Storage while away',
    sitStartDate: '2022-09-15',
    sitEndDate: '2022-09-20',
    ...fieldOverrides,
  };

  return createCompleteMovingExpense({ serviceMemberId, creationDate }, fullFieldOverrides);
};

export { createBaseMovingExpense, createCompleteMovingExpense, createCompleteSITMovingExpense };
