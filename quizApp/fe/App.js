import React, { useState, useEffect } from 'react';
import {
  StyleSheet,
  Text,
  View,
  TouchableOpacity,
  SafeAreaView,
  ScrollView,
  Alert,
} from 'react-native';

const API_BASE_URL = process.env.NODE_ENV === 'production'
  ? 'https://your-api-domain.com/'
  : 'http://localhost:8080/';

export default function App() {
  const [categories, setCategories] = useState([]);
  const [selectedCategory, setSelectedCategory] = useState(null);
  const [questions, setQuestions] = useState([]);
  const [currentQuestionIndex, setCurrentQuestionIndex] = useState(0);
  const [selectedAnswers, setSelectedAnswers] = useState({});
  const [score, setScore] = useState(0);
  const [showResults, setShowResults] = useState(false);
  const [loading, setLoading] = useState(true);
  const [currentView, setCurrentView] = useState('categories');
  const [answerFeedback, setAnswerFeedback] = useState({});
  const [showAnswerFeedback, setShowAnswerFeedback] = useState(false);
  const [buttonEnabled, setButtonEnabled] = useState(false);
  const [runningScore, setRunningScore] = useState(0);

  useEffect(() => {
    fetchCategories();
  }, []);

  const fetchCategories = async () => {
    try {
      setLoading(true);
      console.log('🔍 Attempting to fetch from:', `${API_BASE_URL}categories      // Update the Quiz Screen section (around line 150-200)

      // Quiz Screen
      if (currentView === 'quiz' && questions.length > 0) {
        const currentQuestion = questions[currentQuestionIndex];
        const selectedAnswerId = selectedAnswers[currentQuestionIndex];

        // Add debug logging
        console.log('🎯 Current question data:', currentQuestion);
        console.log('🔤 Available fields:', Object.keys(currentQuestion));

        return (
          <SafeAreaView style={styles.container}>
            <View style={styles.header}>
              <TouchableOpacity onPress={selectedCategory ? backToSubcategories : backToCategories}>
                <Text style={styles.backButton}>← Back</Text>
              </TouchableOpacity>
              <Text style={styles.headerTitle}>
                {selectedCategory ? selectedCategory.category : 'Quiz'}
              </Text>
            </View>

            <View style={styles.questionHeader}>
              <Text style={styles.questionCounter}>
                Question {currentQuestionIndex + 1} of {questions.length}
              </Text>
              <Text style={styles.category}>{currentQuestion.subcategory}</Text>
            </View>

            <ScrollView style={styles.content}>
              {/* FIXED: Use questionText and add fallbacks */}
              <Text style={styles.questionText}>
                {currentQuestion.questionText || currentQuestion.text || currentQuestion.Text || 'Question text not available'}
              </Text>

              {/* Debug info - remove this later */}
              <Text style={{ fontSize: 12, color: 'gray', marginBottom: 10 }}>
                Debug: {JSON.stringify(Object.keys(currentQuestion))}
              </Text>

              <View style={styles.answersContainer}>
                {currentQuestion.answers?.map((answer) => (
                  <TouchableOpacity
                    key={answer.ID}
                    style={[
                      styles.answerButton,
                      selectedAnswerId === answer.ID && styles.selectedAnswer,
                    ]}
                    onPress={() => handleAnswerSelect(answer.ID)}
                  >
                    <Text
                      style={[
                        styles.answerText,
                        selectedAnswerId === answer.ID && styles.selectedAnswerText,
                      ]}
                    >
                      {answer.text}
                    </Text>
                  </TouchableOpacity>
                ))}
              </View>
            </ScrollView>

            <View style={styles.footer}>
              <TouchableOpacity
                style={[
                  styles.nextButton,
                  !selectedAnswerId && styles.disabledButton,
                ]}
                onPress={nextQuestion}
                disabled={!selectedAnswerId}
              >
                <Text style={styles.nextButtonText}>
                  {currentQuestionIndex === questions.length - 1 ? 'Finish' : 'Next'}
                </Text>
              </TouchableOpacity>
            </View>
          </SafeAreaView>
        );
      }`);

      const response = await fetch(`${API_BASE_URL}categories/`);

      if (!response.ok) {
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
      }

      const data = await response.json();

      setCategories(data || []);
      console.log('✅ Categories set successfully!');
      setLoading(false);
    } catch (error) {
      console.error('❌ ERROR Details:', {
        message: error.message,
        name: error.name,
        stack: error.stack
      });
      Alert.alert('Error', `Failed to load categories: ${error.message}`);
      setLoading(false);
    }
  };

  // Add a debug button temporarily
  const DebugButton = () => (
    <TouchableOpacity
      style={[styles.primaryButton, { backgroundColor: '#ff6b6b', marginTop: 10 }]}
      onPress={() => {
        console.log('🐛 Debug button pressed');
        console.log('🌐 API URL:', API_BASE_URL);
        console.log('📊 Current categories:', categories);
        fetchCategories();
      }}
    >
      <Text style={styles.primaryButtonText}>🐛 Debug Fetch</Text>
    </TouchableOpacity>
  );

  const fetchQuestionsByCategory = async (category, subcategory = null) => {
    try {
      setLoading(true);
      let url = `${API_BASE_URL}questions/?category=${encodeURIComponent(category)}`;
      if (subcategory) {
        url += `&subcategory=${encodeURIComponent(subcategory)}`;
      }

      const response = await fetch(url);
      const data = await response.json();
      setQuestions(data || []);
      setCurrentQuestionIndex(0);
      setSelectedAnswers({});
      setScore(0);
      setShowResults(false);
      setCurrentView('quiz');
      setAnswerFeedback({});
      setShowAnswerFeedback(false);
      setButtonEnabled(false);
      setRunningScore(0);
      setLoading(false);
    } catch (error) {
      console.error('Error fetching questions:', error);
      Alert.alert('Error', 'Failed to load questions');
      setLoading(false);
    }
  };

  const handleCategorySelect = (category) => {
    if (category.subcategories && category.subcategories.length > 0) {
      setSelectedCategory(category);
      setCurrentView('subcategories');
    } else {
      fetchQuestionsByCategory(category.category);
    }
  };

  const handleSubcategorySelect = (subcategory) => {
    fetchQuestionsByCategory(selectedCategory.category, subcategory.subcategoryname);
  };

  const handleAnswerSelect = async (answerId) => {
    const currentQuestion = questions[currentQuestionIndex];
    const isMultipleChoice = currentQuestion.type === 'multichoice';

    console.log('🎯 Answer selected:', {
      answerId,
      questionType: currentQuestion.type,
      isMultipleChoice,
      questionText: currentQuestion.questionText
    });

    if (isMultipleChoice) {
      // For multiple choice, allow multiple selections
      const currentSelections = selectedAnswers[currentQuestionIndex] || [];
      // Ensure currentSelections is always an array for multiple choice
      const selectionsArray = Array.isArray(currentSelections) ? currentSelections : [];
      const isAlreadySelected = selectionsArray.includes(answerId);

      let newSelections;
      if (isAlreadySelected) {
        // Remove if already selected
        newSelections = selectionsArray.filter(id => id !== answerId);
      } else {
        // Add to selections
        newSelections = [...selectionsArray, answerId];
      }

      console.log('🔄 Multiple choice selection:', {
        questionIndex: currentQuestionIndex,
        answerId,
        previousSelections: selectionsArray,
        newSelections,
        isAlreadySelected
      });

      setSelectedAnswers({
        ...selectedAnswers,
        [currentQuestionIndex]: newSelections,
      });

      // Enable button immediately for multiple choice (no feedback delay)
      setButtonEnabled(true);
    } else {
      // For single answer questions (true/false, etc.)
      setSelectedAnswers({
        ...selectedAnswers,
        [currentQuestionIndex]: answerId,
      });

      // Disable the button initially to give time to see feedback
      setButtonEnabled(false);

      // Check answer immediately for single-answer questions
      await checkSingleAnswer(answerId);
    }
  };

  const checkSingleAnswer = async (answerId) => {
    // Check if answer is correct using the dedicated answer endpoint
    try {
      const currentQuestion = questions[currentQuestionIndex];

      if (!currentQuestion || !currentQuestion.ID) {
        console.error('❌ Current question is invalid:', currentQuestion);
        return;
      }

      console.log('🚀 Checking answer for question ID:', currentQuestion.ID, 'with answer ID:', answerId);
      const response = await fetch(`${API_BASE_URL}questions/answer/${currentQuestion.ID}/${answerId}`);

      if (!response.ok) {
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
      }

      const answerResult = await response.json();

      console.log('📄 Answer check result:', {
        response: answerResult,
        isCorrect: answerResult.iscorrect || answerResult.IsCorrect,
        correctAnswer: answerResult.correctanswer || answerResult.CorrectAnswer
      });

      const isCorrect = answerResult.iscorrect || answerResult.IsCorrect || false;
      const correctAnswers = answerResult.correctanswer || answerResult.CorrectAnswer || [];
      const correctAnswerId = correctAnswers.length > 0 ? correctAnswers[0].ID : null;

      console.log('🔍 Debug answer checking:', {
        questionId: currentQuestion.ID,
        questionType: currentQuestion.type,
        questionText: currentQuestion.questionText,
        selectedAnswerId: answerId,
        isCorrect,
        correctAnswerId,
        correctAnswers
      });

      setAnswerFeedback({
        ...answerFeedback,
        [currentQuestionIndex]: {
          isCorrect,
          correctAnswerId,
          selectedAnswerId: answerId
        }
      });

      // Update running score for single-answer questions
      if (isCorrect) {
        setRunningScore(prevScore => prevScore + 1);
      }

      setShowAnswerFeedback(true);

      // Enable the button after seeing feedback for 1.5 seconds
      setTimeout(() => {
        setButtonEnabled(true);
      }, 1500);

      // Hide feedback after 4 seconds (giving more time to see the result)
      setTimeout(() => {
        setShowAnswerFeedback(false);
      }, 4000);

    } catch (error) {
      console.error('Error checking answer:', error);
    }
  };

  const nextQuestion = async () => {
    console.log('🚀 Next button clicked:', {
      currentQuestionIndex,
      totalQuestions: questions.length,
      selectedAnswers: selectedAnswers[currentQuestionIndex],
      isLastQuestion: currentQuestionIndex >= questions.length - 1
    });

    // Check and update running score for multiple choice questions
    const currentQuestion = questions[currentQuestionIndex];
    const isMultipleChoice = currentQuestion.type === 'multichoice';

    if (isMultipleChoice) {
      await updateRunningScoreForMultipleChoice();
    }

    // Always hide the answer feedback when moving to next question or finishing
    setShowAnswerFeedback(false);
    setButtonEnabled(false);

    if (currentQuestionIndex < questions.length - 1) {
      setCurrentQuestionIndex(currentQuestionIndex + 1);
    } else {
      calculateScore();
    }
  };

  const updateRunningScoreForMultipleChoice = async () => {
    const currentQuestion = questions[currentQuestionIndex];
    const selectedAnswerData = selectedAnswers[currentQuestionIndex];

    if (!selectedAnswerData || !Array.isArray(selectedAnswerData) || selectedAnswerData.length === 0) {
      return; // No answers selected
    }

    try {
      // Convert array to string format expected by backend
      const answerString = selectedAnswerData.sort().join('');
      const response = await fetch(`${API_BASE_URL}questions/answer/${currentQuestion.ID}/${answerString}`);

      if (response.ok) {
        const answerResult = await response.json();
        const isCorrect = answerResult.iscorrect || answerResult.IsCorrect || false;

        console.log('🏃‍♂️ Running score update for multichoice:', {
          questionId: currentQuestion.ID,
          selectedAnswers: selectedAnswerData,
          isCorrect,
          currentRunningScore: runningScore
        });

        if (isCorrect) {
          setRunningScore(prevScore => prevScore + 1);
        }
      }
    } catch (error) {
      console.error('Error updating running score for multiple choice:', error);
    }
  };

  const calculateScore = async () => {
    let correctAnswers = 0;

    for (let i = 0; i < questions.length; i++) {
      const question = questions[i];
      const selectedAnswerData = selectedAnswers[i];
      const isMultipleChoice = question.type === 'multichoice';

      if (isMultipleChoice && Array.isArray(selectedAnswerData)) {
        // For multiple choice, check if all selected answers are correct
        // This uses a more complex scoring logic for multiple selections
        try {
          // Convert array to string format expected by backend
          const answerString = selectedAnswerData.sort().join('');
          const response = await fetch(`${API_BASE_URL}questions/answer/${question.ID}/${answerString}`);

          if (!response.ok) {
            console.error(`Failed to check multiple choice answer for question ${question.ID}: ${response.status} ${response.statusText}`);
            continue;
          }

          const answerResult = await response.json();
          const isCorrectForScore = answerResult.iscorrect || answerResult.IsCorrect || false;

          console.log(`📊 Multiple choice score Q${i + 1}:`, {
            questionId: question.ID,
            selectedAnswers: selectedAnswerData,
            answerString,
            isCorrectForScore,
            answerResult
          });

          if (isCorrectForScore) {
            correctAnswers++;
          }
        } catch (error) {
          console.error('Error checking multiple choice answer:', error);
        }
      } else {
        // Single answer question
        const selectedAnswerId = Array.isArray(selectedAnswerData) ? selectedAnswerData[0] : selectedAnswerData;

        if (!selectedAnswerId) continue;

        try {
          const response = await fetch(`${API_BASE_URL}questions/answer/${question.ID}/${selectedAnswerId}`);

          if (!response.ok) {
            console.error(`Failed to check answer for question ${question.ID}: ${response.status} ${response.statusText}`);
            continue;
          }

          const answerResult = await response.json();
          const isCorrectForScore = answerResult.iscorrect || answerResult.IsCorrect || false;

          console.log(`📊 Single answer score Q${i + 1}:`, {
            questionId: question.ID,
            selectedAnswerId,
            isCorrectForScore,
            answerResult
          });

          if (isCorrectForScore) {
            correctAnswers++;
          }
        } catch (error) {
          console.error('Error checking single answer:', error);
        }
      }
    }

    setScore(correctAnswers);
    setShowResults(true);
  };

  const resetToCategories = () => {
    setCurrentView('categories');
    setSelectedCategory(null);
    setQuestions([]);
    setCurrentQuestionIndex(0);
    setSelectedAnswers({});
    setScore(0);
    setShowResults(false);
    setAnswerFeedback({});
    setShowAnswerFeedback(false);
    setButtonEnabled(false);
    setRunningScore(0);
  };

  const backToCategories = () => {
    setCurrentView('categories');
    setSelectedCategory(null);
  };

  const backToSubcategories = () => {
    setCurrentView('subcategories');
  };

  // Loading Screen
  if (loading) {
    return (
      <SafeAreaView style={styles.container}>
        <View style={styles.centerContainer}>
          <Text style={styles.loadingText}>Loading...</Text>
        </View>
      </SafeAreaView>
    );
  }

  // Results Screen
  if (showResults) {
    return (
      <SafeAreaView style={styles.container}>
        <View style={styles.header}>
          <TouchableOpacity onPress={resetToCategories}>
            <Text style={styles.backButton}>← Back to Categories</Text>
          </TouchableOpacity>
          <Text style={styles.headerTitle}>Quiz Results</Text>
        </View>
        <View style={styles.resultsContainer}>
          <Text style={styles.resultsTitle}>Quiz Complete!</Text>
          <Text style={styles.scoreText}>
            Your Score: {score} / {questions.length}
          </Text>
          <Text style={styles.percentageText}>
            {Math.round((score / questions.length) * 100)}%
          </Text>
          <TouchableOpacity style={styles.primaryButton} onPress={resetToCategories}>
            <Text style={styles.primaryButtonText}>Back to Categories</Text>
          </TouchableOpacity>
        </View>
      </SafeAreaView>
    );
  }

  // Quiz Screen
  if (currentView === 'quiz' && questions.length > 0) {
    const currentQuestion = questions[currentQuestionIndex];
    const selectedAnswerData = selectedAnswers[currentQuestionIndex];
    const isMultipleChoice = currentQuestion.type === 'multichoice';

    // Handle both single and multiple selections
    const selectedAnswerIds = isMultipleChoice && Array.isArray(selectedAnswerData) ?
      selectedAnswerData :
      (selectedAnswerData ? [selectedAnswerData] : []);

    return (
      <SafeAreaView style={styles.container}>
        <View style={styles.header}>
          <TouchableOpacity onPress={selectedCategory ? backToSubcategories : backToCategories}>
            <Text style={styles.backButton}>← Back</Text>
          </TouchableOpacity>
          <Text style={styles.headerTitle}>
            {selectedCategory ? selectedCategory.category : 'Quiz'}
          </Text>
        </View>

        <View style={styles.questionHeader}>
          <Text style={styles.questionCounter}>
            Question {currentQuestionIndex + 1} of {questions.length}
          </Text>
          <Text style={styles.category}>{currentQuestion.subcategory}</Text>
          <View style={styles.scoreContainer}>
            <Text style={styles.runningScore}>
              Score: {runningScore} / {questions.length}
            </Text>
            <Text style={styles.percentageScore}>
              {questions.length > 0 ? Math.round((runningScore / questions.length) * 100) : 0}%
            </Text>
          </View>
        </View>

        <ScrollView style={styles.content}>
          <Text style={styles.questionText}>{currentQuestion.questionText}</Text>

          {/* Answer feedback display */}
          {showAnswerFeedback && answerFeedback[currentQuestionIndex] && (
            <View style={[
              styles.feedbackContainer,
              answerFeedback[currentQuestionIndex]?.isCorrect
                ? styles.correctFeedback
                : styles.incorrectFeedback
            ]}>
              <Text style={styles.feedbackText}>
                {answerFeedback[currentQuestionIndex]?.isCorrect
                  ? '✅ Correct!'
                  : '❌ Incorrect'}
              </Text>
            </View>
          )}

          <View style={styles.answersContainer}>
            {isMultipleChoice && (
              <Text style={styles.instructionText}>
                💡 You can select multiple answers for this question
              </Text>
            )}
            {currentQuestion.answers?.map((answer) => {
              const feedback = answerFeedback[currentQuestionIndex];
              const isSelected = selectedAnswerIds.includes(answer.ID);
              const isCorrect = feedback?.correctAnswerId === answer.ID;
              const isIncorrectSelection = feedback && isSelected && !feedback.isCorrect;

              return (
                <TouchableOpacity
                  key={answer.ID}
                  style={[
                    styles.answerButton,
                    isSelected && styles.selectedAnswer,
                    feedback && isCorrect && styles.correctAnswer,
                    feedback && isIncorrectSelection && styles.incorrectAnswer,
                  ]}
                  onPress={() => handleAnswerSelect(answer.ID)}
                >
                  <Text
                    style={[
                      styles.answerText,
                      isSelected && styles.selectedAnswerText,
                      feedback && isCorrect && styles.correctAnswerText,
                      feedback && isIncorrectSelection && styles.incorrectAnswerText,
                    ]}
                  >
                    {answer.text}
                    {feedback && isCorrect && ' ✅'}
                    {feedback && isIncorrectSelection && ' ❌'}
                  </Text>
                </TouchableOpacity>
              );
            })}
          </View>
        </ScrollView>

        <View style={styles.footer}>
          <TouchableOpacity
            style={[
              styles.nextButton,
              (selectedAnswerIds.length === 0 || (!isMultipleChoice && !buttonEnabled)) && styles.disabledButton,
            ]}
            onPress={nextQuestion}
            disabled={selectedAnswerIds.length === 0 || (!isMultipleChoice && !buttonEnabled)}
          >
            <Text style={styles.nextButtonText}>
              {currentQuestionIndex === questions.length - 1 ? 'Finish' : 'Next'}
            </Text>
          </TouchableOpacity>
        </View>
      </SafeAreaView>
    );
  }

  // Subcategories Screen
  if (currentView === 'subcategories' && selectedCategory) {
    return (
      <SafeAreaView style={styles.container}>
        <View style={styles.header}>
          <TouchableOpacity onPress={backToCategories}>
            <Text style={styles.backButton}>← Categories</Text>
          </TouchableOpacity>
          <Text style={styles.headerTitle}>{selectedCategory.category}</Text>
        </View>

        <ScrollView style={styles.content}>
          <Text style={styles.sectionTitle}>Select a Subcategory:</Text>

          <View style={styles.listContainer}>
            {selectedCategory.subcategories.map((subcategory, index) => (
              <TouchableOpacity
                key={index}
                style={styles.categoryCard}
                onPress={() => handleSubcategorySelect(subcategory)}
              >
                <View style={styles.categoryContent}>
                  <Text style={styles.categoryTitle}>
                    {subcategory.subcategoryname}
                  </Text>
                  <Text style={styles.categoryArrow}>→</Text>
                </View>
              </TouchableOpacity>
            ))}
          </View>
        </ScrollView>
      </SafeAreaView>
    );
  }

  // Categories Screen (Default)
  return (
    <SafeAreaView style={styles.container}>
      <View style={styles.header}>
        <Text style={styles.headerTitle}>Quizzie Categories</Text>
      </View>

      <ScrollView style={styles.content}>
        <Text style={styles.sectionTitle}>Select a Category:</Text>

        {categories.length === 0 ? (
          <View style={styles.centerContainer}>
            <Text style={styles.errorText}>No categories available</Text>
            <TouchableOpacity style={styles.primaryButton} onPress={fetchCategories}>
              <Text style={styles.primaryButtonText}>Retry</Text>
            </TouchableOpacity>
            <DebugButton />
          </View>
        ) : (
          <View style={styles.listContainer}>
            {categories.map((category, index) => (
              <TouchableOpacity
                key={index}
                style={styles.categoryCard}
                onPress={() => handleCategorySelect(category)}
              >
                <View style={styles.categoryContent}>
                  <View>
                    <Text style={styles.categoryTitle}>{category.category}</Text>
                    <Text style={styles.subcategoryCount}>
                      {category.subcategories?.length || 0} subcategories
                    </Text>
                  </View>
                  <Text style={styles.categoryArrow}>→</Text>
                </View>
              </TouchableOpacity>
            ))}
          </View>
        )}
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f8f9fa',
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    padding: 20,
    backgroundColor: '#fff',
    borderBottomWidth: 1,
    borderBottomColor: '#e9ecef',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 3,
    elevation: 3,
  },
  backButton: {
    fontSize: 16,
    color: '#007AFF',
    fontWeight: '600',
  },
  headerTitle: {
    fontSize: 20,
    fontWeight: 'bold',
    color: '#2c3e50',
    flex: 1,
    textAlign: 'center',
    marginLeft: -50,
  },
  content: {
    flex: 1,
    padding: 20,
  },
  sectionTitle: {
    fontSize: 24,
    fontWeight: 'bold',
    color: '#2c3e50',
    marginBottom: 20,
  },
  listContainer: {
    gap: 15,
  },
  categoryCard: {
    backgroundColor: '#fff',
    borderRadius: 12,
    padding: 20,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 8,
    elevation: 3,
  },
  categoryContent: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  categoryTitle: {
    fontSize: 18,
    fontWeight: '600',
    color: '#2c3e50',
  },
  subcategoryCount: {
    fontSize: 14,
    color: '#6c757d',
    marginTop: 4,
  },
  categoryArrow: {
    fontSize: 20,
    color: '#007AFF',
    fontWeight: 'bold',
  },
  centerContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
  loadingText: {
    fontSize: 18,
    color: '#6c757d',
  },
  errorText: {
    fontSize: 18,
    textAlign: 'center',
    color: '#dc3545',
    marginBottom: 20,
  },
  primaryButton: {
    backgroundColor: '#007AFF',
    padding: 15,
    borderRadius: 10,
    minWidth: 120,
    alignItems: 'center',
  },
  primaryButtonText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: '600',
  },
  questionHeader: {
    padding: 20,
    backgroundColor: '#fff',
    borderBottomWidth: 1,
    borderBottomColor: '#e9ecef',
  },
  questionCounter: {
    fontSize: 14,
    color: '#6c757d',
    textAlign: 'center',
  },
  category: {
    fontSize: 16,
    fontWeight: 'bold',
    color: '#495057',
    textAlign: 'center',
    marginTop: 5,
  },
  questionText: {
    fontSize: 18,
    fontWeight: '600',
    color: '#2c3e50',
    marginBottom: 30,
    lineHeight: 26,
  },
  answersContainer: {
    gap: 15,
  },
  answerButton: {
    backgroundColor: '#fff',
    padding: 15,
    borderRadius: 10,
    borderWidth: 2,
    borderColor: '#e9ecef',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 1 },
    shadowOpacity: 0.1,
    shadowRadius: 3,
    elevation: 2,
  },
  selectedAnswer: {
    borderColor: '#007AFF',
    backgroundColor: '#f0f8ff',
  },
  answerText: {
    fontSize: 16,
    color: '#495057',
  },
  selectedAnswerText: {
    color: '#007AFF',
    fontWeight: '600',
  },
  footer: {
    padding: 20,
    backgroundColor: '#fff',
    borderTopWidth: 1,
    borderTopColor: '#e9ecef',
  },
  nextButton: {
    backgroundColor: '#28a745',
    padding: 15,
    borderRadius: 10,
    alignItems: 'center',
  },
  disabledButton: {
    backgroundColor: '#6c757d',
  },
  nextButtonText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: '600',
  },
  resultsContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 20,
  },
  resultsTitle: {
    fontSize: 28,
    fontWeight: 'bold',
    color: '#2c3e50',
    marginBottom: 20,
  },
  scoreText: {
    fontSize: 24,
    color: '#495057',
    marginBottom: 10,
  },
  percentageText: {
    fontSize: 48,
    fontWeight: 'bold',
    color: '#28a745',
    marginBottom: 30,
  },
  totalQuestions: {
    fontSize: 12,
    color: '#6c757d',
    textAlign: 'center',
    marginTop: 5,
    fontStyle: 'italic',
  },
  feedbackContainer: {
    padding: 15,
    borderRadius: 10,
    marginBottom: 20,
    alignItems: 'center',
  },
  correctFeedback: {
    backgroundColor: '#d4edda',
    borderColor: '#28a745',
    borderWidth: 1,
  },
  incorrectFeedback: {
    backgroundColor: '#f8d7da',
    borderColor: '#dc3545',
    borderWidth: 1,
  },
  feedbackText: {
    fontSize: 16,
    fontWeight: 'bold',
  },
  correctAnswer: {
    backgroundColor: '#d4edda',
    borderColor: '#28a745',
  },
  incorrectAnswer: {
    backgroundColor: '#f8d7da',
    borderColor: '#dc3545',
  },
  correctAnswerText: {
    color: '#155724',
    fontWeight: 'bold',
  },
  incorrectAnswerText: {
    color: '#721c24',
    fontWeight: 'bold',
  },
  instructionText: {
    fontSize: 14,
    color: '#6c757d',
    fontStyle: 'italic',
    textAlign: 'center',
    marginBottom: 15,
    backgroundColor: '#f8f9fa',
    padding: 10,
    borderRadius: 8,
    borderLeftWidth: 4,
    borderLeftColor: '#007AFF',
  },
  scoreContainer: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginTop: 10,
    paddingTop: 10,
    borderTopWidth: 1,
    borderTopColor: '#e9ecef',
  },
  runningScore: {
    fontSize: 16,
    fontWeight: 'bold',
    color: '#28a745',
  },
  percentageScore: {
    fontSize: 14,
    fontWeight: '600',
    color: '#007AFF',
  },
});
