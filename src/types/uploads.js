import * as Yup from 'yup';

const uploadShape = Yup.object({
  id: Yup.string(),
  created_at: Yup.string(),
  bytes: Yup.number(),
  url: Yup.string(),
  filename: Yup.string(),
  content_type: Yup.string(),
});

export default uploadShape;
