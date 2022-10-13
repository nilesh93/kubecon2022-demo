import express from 'express'
import dotenv from 'dotenv'
import { taskGroupRoute } from './taskgroup'

dotenv.config();
const app = express();

app.use(express.json());
//app.use(`/task`, taskRoute);
app.use(`/taskgroup`, taskGroupRoute);
//app.use('/taskstatus', taskStatusRoute);

const server = app.listen(8080, () =>
  console.log(`Server ready at: http://localhost:8080`),
)
