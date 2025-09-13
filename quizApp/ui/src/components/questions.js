import React, {useState, useEffect } from 'react';
import AnswerCheck from './answerChk';
import NextButton from './nextButton';
import SubmitButton from './submitButton';

export default function Questions(props) {
    const [QuestionList, setQuestions] = useState([]);
    const [NumQuestions, setNumQuestions] = useState(0);
    const [QuestionIdx, setQuestionIdx] = useState(0);

    const [PostResponse, setPostResponse ] = useState([]);
    const [AnswerResp, setAnswerResp] = useState([])

    let questionUrl =  `http://localhost:5000/questions/${props.catName}/${props.subCatName}`
    let answerURL = "http://localhost:5000/questions/answer/"

    //Next Question
    function nextQuestion() {
      let newQuestionidx = QuestionIdx + 1
      setQuestionIdx(newQuestionidx)
      setAnswerResp([])
      setPostResponse([])
    }

    //CheckBox Checked
    const handleChange = (e) => {
      const { value, checked } = e.target;
      console.log(`${value} is ${checked}`);

      if (checked) {
        setPostResponse([...PostResponse, value]);
      }

      else {
        setPostResponse(
          PostResponse.filter(
                  (e) => e !== value
              )
          );
      }
    }

    const fetchQuestions = async () => {
        try {
          const response = await fetch(questionUrl, {
            method: 'GET'
          });
          if (!response.ok) {
            throw new Error('Network response was not ok');
          }

          const data = await response.json();
          setQuestions(data);
          let numQuestions = data.length
          setNumQuestions(numQuestions)
        } catch (error) {
          console.error('Error:', error);
        }
      }

    const submitAnswer = async (id) => {
      let answer = PostResponse.join('');
      let AnswerURL = `${answerURL}${id}/${answer}`
      try {
        const response = await fetch(AnswerURL, {
          method: 'GET'
        });
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
        const data = await response.json();
        setAnswerResp([data]);
      } catch (error) {
        console.error('Error:', error);
      }
//      setCorrectAnswerCount((prevValue) => prevValue + 1)
    }

    const displayAnswers = AnswerResp && AnswerResp.map((answer,idx) =>

   <div>
       { answer.iscorrect &&
        <div key={idx} className='correctAnswer'>
          <h6><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-check" viewBox="0 0 16 16">
  <path d="M10.97 4.97a.75.75 0 0 1 1.07 1.05l-3.99 4.99a.75.75 0 0 1-1.08.02L4.324 8.384a.75.75 0 1 1 1.06-1.06l2.094 2.093 3.473-4.425a.267.267 0 0 1 .02-.022z"/>
            </svg>
          Correct! You submitted the correct answer</h6>
        </div>
      }
      { !answer.iscorrect &&
        <div key={idx} className='incorrectAnswer'>
        <h6>
        Sorry, that is an <strong>incorrect</strong> Answer.</h6>
        <div>
          <div>
          <h7>Actual Answers</h7>
            {answer.correctanswer.map((answer,idx) =>
            <div>
              <p>- {answer.text}</p>
            </div>
            )}
        </div>
      </div>
      </div>
      }
      </div>
    );

    const displayQuestion = QuestionList.map((question,idx) =>
    <div>
    <div className="question">
        <p>{question.questionText}</p>
        <hr ></hr>
        {question.answers.map((answer,idx) =>
        <div>
          <AnswerCheck
            answer={answer}
            idx={idx}
            id={answer.ID}
            handleChange={handleChange}
          />
        </div>
        )}
    </div>
    <div className='answer'>
        {displayAnswers}
    </div>
       <SubmitButton
          submitAnswer={submitAnswer}
          id={question.ID}
       />
      <NextButton
        NextQuestion={nextQuestion}
        NumQuestions={NumQuestions}
        QuestionIdx={QuestionIdx}
      />
       <p className='id'>Question Id: {question.ID}</p>
    </div>
    )[QuestionIdx];

    useEffect(() => {
      fetchQuestions()
       }, [AnswerResp]);

    return (
    <div>
      <div className='row page-header'>
        <div className='col-md'> <h2 id="list-heading"><a className="parentCat" href='/'>{props.catName} &#92;</a> {props.subCatName}</h2></div>
        <div className='col-md'><p>Question: {QuestionIdx + 1} of {NumQuestions}</p></div>
      </div>

        <hr></hr>
        {displayQuestion}
      </div>
    );
  }
